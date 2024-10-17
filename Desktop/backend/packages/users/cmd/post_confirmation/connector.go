//go:build !test
// +build !test

package main

import (
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
)

func connector(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	logger := logging.Must(logging.NewLoggerFromEvent(event))
	defer logging.SyncLogger(logger)

	logger.Info("post_confirmation started")

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	idp := cognitoidentityprovider.New(sess)

	sErr := handler(HandlerInput{
		UserName:     event.UserName,
		PoolID:       event.UserPoolID,
		Idp:          idp,
		Logger:       logger,
		DefaultGroup: os.Getenv("DEFAULT_GROUP"),
	})
	if sErr != nil {
		return events.CognitoEventUserPoolsPostConfirmation{}, errors.New(sErr.Error())
	}

	logger.Info("post_confirmation completed")

	return event, nil
}

func main() {
	lambda.Start(connector)
}
