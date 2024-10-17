//go:build !test
// +build !test

package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"go.uber.org/zap"
)

func connector(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	logger := logging.Must(logging.NewLoggerFromEvent(event))
	defer logging.SyncLogger(logger)

	logger = logger.With(zap.String("poolID", event.UserPoolID))

	logger.Info("pre_sign_up started")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)
	if err != nil {
		return event, err
	}

	organizationName := event.Request.UserAttributes["custom:organization"]

	sErr := handler(HandlerInput{
		UserName:         event.UserName,
		PoolID:           event.UserPoolID,
		DB:               db,
		Logger:           logger,
		OrganizationName: organizationName,
	})
	if sErr != nil {
		return events.CognitoEventUserPoolsPreSignup{}, errors.New(sErr.Error())
	}

	logger.Info("pre_sign_up completed")

	event.Response.AutoVerifyEmail = true
	event.Response.AutoConfirmUser = true
	return event, nil
}

func main() {
	lambda.Start(connector)
}
