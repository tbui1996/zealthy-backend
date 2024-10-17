package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/request"
	"go.uber.org/zap"
)

type ServiceReceiveRequest struct {
	Name             string
	ReceiveQueueName string
	SendQueueName    string
	SQS              sqsiface.SQSAPI
	Event            events.APIGatewayWebsocketProxyRequest
	Logger           *zap.Logger
}

func Handler(req ServiceReceiveRequest) error {
	var requestPayload request.Payload
	err := json.Unmarshal([]byte(req.Event.Body), &requestPayload)
	if err != nil {
		req.Logger.Error("couldn't unmarshal received message: " + req.Event.Body)
	}

	req.Logger.Debug("processing payload: " + requestPayload.Payload)

	queue, err := req.SQS.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(req.ReceiveQueueName),
	})

	if err != nil {
		return err
	}

	logger := logging.LoggerFieldsFromEvent(req.Event)
	if logger.Error != nil {
		return logger.Error
	}

	loggerFieldsBytes, err := json.Marshal(logger.Fields)
	if err != nil {
		return err
	}

	req.Logger.Debug("sending the following logger fields" + string(loggerFieldsBytes))

	_, err = req.SQS.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Source": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("sonar-service-router-%s", req.Name)),
			},
			"ConnectionId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(req.Event.RequestContext.ConnectionID),
			},
			"LoggerFields": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(loggerFieldsBytes)),
			},
		},
		MessageBody: aws.String(requestPayload.Payload),
		QueueUrl:    queue.QueueUrl,
	})

	if err != nil {
		return err
	}

	return nil
}
