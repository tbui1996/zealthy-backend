package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/request"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ReceiveHandlerTestSuite struct {
	suite.Suite
}

func getItemOutput() *dynamodb.GetItemOutput {
	return &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String("1"),
			},
			"UserID": {
				S: aws.String("1"),
			},
		},
	}
}

func getQueryOutput() *dynamodb.QueryOutput {
	itemOne := map[string]*dynamodb.AttributeValue{
		"UserID": {
			S: aws.String("1"),
		},
		"CreatedTimestamp": {
			N: aws.String("1"),
		},
		"DeleteTimestamp": {
			N: aws.String("1"),
		},
		"Message": {
			S: aws.String("Hello"),
		},
	}

	return &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{itemOne},
	}
}

func getEvent() *events.SQSEvent {
	req := request.RouterTypeRequest{
		Type:    "undelivered_messages",
		Message: "test",
	}

	b, _ := json.Marshal(req)

	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"ConnectionId": {StringValue: aws.String("1")},
				},
				Body: string(b),
			},
		},
	}
}

func getBadEventType() *events.SQSEvent {
	req := request.RouterTypeRequest{
		Type:    "test",
		Message: "test",
	}

	b, _ := json.Marshal(req)

	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"ConnectionId": {StringValue: aws.String("1")},
				},
				Body: string(b),
			},
		},
	}
}

func getBadEventNoConnection() *events.SQSEvent {
	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"ConnectionId": {StringValue: nil},
				},
			},
		},
	}
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_Success() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: aws.String("url"),
	}, nil)

	/* Undelivered Handler */
	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(getQueryOutput(), nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	mockDB.On("DeleteItem", mock.Anything).Return(nil, nil)
	/* Undelivered Handler */

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, nil)

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getEvent(),
	})

	suite.Empty(errorArray)
	suite.Nil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_FailOnGetQueueUrl() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getEvent(),
	})

	suite.Empty(errorArray)
	suite.NotNil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_FailOnNoConnection() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: aws.String("url"),
	}, nil)

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getBadEventNoConnection(),
	})

	suite.NotEmpty(errorArray)
	suite.Nil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_FailOnInvalidType() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: aws.String("url"),
	}, nil)

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getBadEventType(),
	})

	suite.NotEmpty(errorArray)
	suite.Nil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_FailOnHandleUndelivered() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: aws.String("url"),
	}, nil)

	/* Undelivered Handler */
	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), errors.New("FAKE TEST ERROR, IGNORE"))
	/* Undelivered Handler */

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, nil)

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getEvent(),
	})

	suite.NotEmpty(errorArray)
	suite.Nil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ReceiveHandlerTestSuite) TestReceiveHandler_FailOnDeleteMessage() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: aws.String("url"),
	}, nil)

	/* Undelivered Handler */
	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(getQueryOutput(), nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	mockDB.On("DeleteItem", mock.Anything).Return(nil, nil)
	/* Undelivered Handler */

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	errorArray, err := Handler(ReceiveRequest{
		Dynamo: mockDB,
		API:    mockApi,
		SQS:    mockSqs,
		Event:  *getEvent(),
	})

	suite.NotEmpty(errorArray)
	suite.Nil(err)

	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertCalled(suite.T(), "DeleteItem", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func TestReceiveHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiveHandlerTestSuite))
}
