//go:build !test
// +build !test

package main

import (
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"context"
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	config, err := requestConfig.Must(requestConfig.NewRequestConfigNoLogger(ctx, sqsEvent)).ToSQSEventRequest()

	if err != nil {
		return fmt.Errorf("unable to get config for session (%s)", err.Error())
	}

	sqsSvc := sqs.New(config.Session)

	dynamoSvc := dynamodb.New(config.Session)

	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String(os.Getenv("WEBSOCKET_URL")),
	})

	errorArray, err := Handler(ReceiveRequest{
		SQS:    sqsSvc,
		Dynamo: dynamoSvc,
		API:    api,
		Event:  sqsEvent,
	})

	if err != nil {
		return err
	}

	if len(errorArray) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errorArray, "\n"))
}

func main() {
	lambda.Start(HandleRequest)
}
