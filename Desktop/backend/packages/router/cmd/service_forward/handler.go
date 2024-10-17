package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/forward"
	"go.uber.org/zap"
)

type HandleMessageInput struct {
	Message                 events.SQSMessage
	SQS                     sqsiface.SQSAPI
	Name                    string
	DynamoDB                dynamodbiface.DynamoDBAPI
	ApiGatewayManagementApi apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
	ReceiveQueueName        string
	SendQueueName           string
	Logger                  *zap.Logger
}

func Handler(input *HandleMessageInput) error {
	forwarder, err := forward.NewForwarderFromSQS(forward.ForwarderSqsDTO{
		Message:                 input.Message,
		DynamoDB:                input.DynamoDB,
		ApiGatewayManagementApi: input.ApiGatewayManagementApi,
		Logger:                  input.Logger,
	})

	if err != nil {
		return err
	}

	input.Logger = input.Logger.With(
		zap.String("name", input.Name),
		zap.String("messageID", input.Message.MessageId),
	)

	input.Logger.Info("processing forward message")

	err = forwarder.Forward()
	if err != nil {
		return err
	}

	input.Logger.Info("forwarded message to recipients")

	queue, err := input.SQS.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(input.SendQueueName),
	})

	if err != nil {
		return err
	}

	_, err = input.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queue.QueueUrl,
		ReceiptHandle: aws.String(input.Message.ReceiptHandle),
	})

	if err != nil {
		input.Logger.Error("deleting message from send queue: " + err.Error())
		return err
	}

	return nil
}
