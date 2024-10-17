//go:build !test
// +build !test

package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
)

func HandleRequest(event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	logger := logging.Must(logging.NewLoggerFromEvent(event))
	defer func() {
		if err := logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	logger.Info("external_oh_authorizer started")

	svc := authorizer.Jwt{}

	input := HandleRequestInput{
		Event:  event,
		Jwt:    &svc,
		Logger: logger,
	}

	resp, sErr := handleRequest(input)
	if sErr != nil {
		logger.Error(sErr.Error())
	}

	logger.Info("external_oh_authorizer completed, sending response")

	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
