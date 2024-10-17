//go:build !test
// +build !test

package main

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"net/http"
	"os"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"context"
	"log"
)

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayProxyResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	name := os.Getenv("CONTEXT")
	config.Logger = config.Logger.With(zap.String("context", name))
	config.Logger.Info("service_receive called")

	sqsClient := sqs.New(config.Session)
	err = Handler(ServiceReceiveRequest{
		Name:             name,
		ReceiveQueueName: os.Getenv("RECEIVE_QUEUE"),
		SendQueueName:    os.Getenv("SEND_QUEUE"),
		SQS:              sqsClient,
		Event:            config.Event,
		Logger:           config.Logger,
	})

	if err != nil {
		config.Logger.Error("Error receiving payload: " + err.Error())
		return exception.NewSonarError(http.StatusInternalServerError, "Error receiving "+name+" payload!").ToAPIGatewayProxyResponse(), nil
	}

	config.Logger.Info("service_receive completed, sending response")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
