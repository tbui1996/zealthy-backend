//go:build !test
// +build !test

package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"go.uber.org/zap"
)

func connector(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	sess := session.Must(session.NewSession())
	idp := cognitoidentityprovider.New(sess)

	environment := os.Getenv("ENVIRONMENT")
	poolID := event.UserPoolID
	username := event.UserName

	logger := logging.Must(logging.NewBasicLogger())
	defer logging.SyncLogger(logger)

	logger.Info("post_authentication called")

	logger = logger.With(zap.String("poolID", poolID))
	logger = logger.With(zap.String("username", username))

	logger.Debug("parsing event")
	pe := parse(event, environment, logger)

	logger.Debug("validating event")
	ve, sErr := validate(ValidateEventInput{
		ParsedEvent: pe,
		Idp:         idp,
		Environment: environment,
		Username:    username,
		PoolID:      poolID,
	})
	if sErr != nil {
		logger.Error("validating request: " + sErr.Error())
		return event, sErr
	}

	logger = logger.With(zap.String("okta_group", ve.OktaGroup))
	logger = logger.With(zap.String("cognito_groups", strings.Join(ve.CognitoGroups, ", ")))
	logger.Debug("handling event")

	err := handler(HandlerInput{
		Idp:            idp,
		ValidatedEvent: ve,
		Username:       username,
		PoolID:         poolID,
		Logger:         logger,
	})
	if err != nil {
		logger.Error("something went wrong handling pre authentication" + err.Error())
		return event, err
	}

	logger.Info("post_authentication completed, sending response")
	return event, nil
}

func main() {
	lambda.Start(connector)
}
