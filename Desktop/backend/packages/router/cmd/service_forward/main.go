//go:build !test
// +build !test

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"os"
	"strings"

	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	config, err := requestConfig.Must(requestConfig.NewRequestConfigNoLogger(ctx, sqsEvent)).ToSQSEventRequest()

	if err != nil {
		return fmt.Errorf("unable to get session config for sqs event %s", err.Error())
	}

	sqsClient := sqs.New(config.Session)
	svc := dynamodb.New(config.Session)

	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String(os.Getenv("WEBSOCKET_URL")),
	})

	// Holds all errors generated while parsing sqs events
	errors := make([]string, 0)
	for _, message := range sqsEvent.Records {
		var loggerFields logging.LoggerFields
		logger, err := loggerFields.FromSQSMessage(message)

		if err != nil {
			errors = append(errors, err.Error())
		}

		if logger == nil {
			continue
		}

		input := &HandleMessageInput{
			Message:                 message,
			SQS:                     sqsClient,
			Name:                    os.Getenv("CONTEXT"),
			DynamoDB:                svc,
			ApiGatewayManagementApi: api,
			ReceiveQueueName:        os.Getenv("RECEIVE_QUEUE"),
			SendQueueName:           os.Getenv("SEND_QUEUE"),
			Logger:                  logger,
		}

		err = Handler(input)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errors, "\n"))
}

func main() {
	lambda.Start(HandleRequest)
}
