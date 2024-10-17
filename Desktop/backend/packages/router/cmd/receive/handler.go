package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/request"
)

type ReceiveRequest struct {
	Event  events.SQSEvent
	Dynamo dynamodbiface.DynamoDBAPI
	API    apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
	SQS    sqsiface.SQSAPI
}

func Handler(req ReceiveRequest) ([]string, error) {
	queue, err := req.SQS.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("sonar-service-router-receive"),
	})

	if err != nil {
		return nil, err
	}

	errs := make([]string, 0)
	for _, message := range req.Event.Records {
		connectionID, ok := message.MessageAttributes["ConnectionId"]
		if !ok || connectionID.StringValue == nil {
			errs = append(errs, "expected connection ID to be in the message attributes")
			continue
		}

		var routerTypeRequest request.RouterTypeRequest
		if err := json.Unmarshal([]byte(message.Body), &routerTypeRequest); err != nil {
			errs = append(errs, err.Error())
			continue
		}

		switch routerTypeRequest.Type {
		case "undelivered_messages":
			err := UndeliveredHandler(req.Dynamo, req.API, *connectionID.StringValue)

			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
		default:
			errs = append(errs, fmt.Sprintf("invalid message type %s", routerTypeRequest.Type))
			continue
		}

		// Delete message after receive
		_, err := req.SQS.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      queue.QueueUrl,
			ReceiptHandle: aws.String(message.ReceiptHandle),
		})

		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}

	return errs, nil
}
