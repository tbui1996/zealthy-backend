//go:build !test
// +build !test

package main

import (
	"context"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	commonEvents "github.com/circulohealth/sonar-backend/packages/common/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayProxyResponse(), nil
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("connect called")

	db := &dynamo.DynamoDatabase{
		TableName: dynamo.SonarInternalWebsocketConnections,
	}

	sErr := handleRequest(HandleRequestInput{
		DB:     db,
		Event:  config.Event,
		Logger: config.Logger,
		EventPublisher: &commonEvents.EventPublisher{
			EventBridge: eventbridge.New(config.Session),
		},
	})

	if sErr != nil {
		config.Logger.Error(sErr.Error())
		return sErr.ToAPIGatewayProxyResponse(), nil
	}

	config.Logger.Info("connect completed, sending response")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
