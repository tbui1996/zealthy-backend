package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceForwardTestSuite struct {
	suite.Suite
}

func sqsMessage() events.SQSMessage {
	action := "action"
	procedure := "procedure"
	source := "action"

	l := logging.LoggerFields{
		SourceIP:  "1",
		RouteKey:  "1",
		RequestID: "1",
		UserAgent: "1",
		UserID:    "1",
		Email:     "1",
	}

	b, _ := json.Marshal(l)

	loggerFields := string(b)

	return events.SQSMessage{MessageAttributes: map[string]events.SQSMessageAttribute{
		"Action":       {StringValue: &action},
		"Procedure":    {StringValue: &procedure},
		"Source":       {StringValue: &source},
		"LoggerFields": {StringValue: &loggerFields},
	}}
}

func connItemOne() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"ConnectionId": {
			S: aws.String("1"),
		},
		"UserID": {
			S: aws.String("1"),
		},
	}
}

func (suite *ServiceForwardTestSuite) TestServiceForward_Success() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{connItemOne()},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	url := "url"
	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: &url,
	}, nil)

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, nil)

	err := Handler(&HandleMessageInput{
		Message:                 sqsMessage(),
		SQS:                     mockSqs,
		Name:                    "test",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		ReceiveQueueName:        "receive",
		SendQueueName:           "send",
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ServiceForwardTestSuite) TestServiceForward_FailOnScan() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockDB.On("Scan", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(&HandleMessageInput{
		Message:                 sqsMessage(),
		SQS:                     mockSqs,
		Name:                    "test",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		ReceiveQueueName:        "receive",
		SendQueueName:           "send",
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ServiceForwardTestSuite) TestServiceForward_FailOnPostToConnection() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{connItemOne()},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(&HandleMessageInput{
		Message:                 sqsMessage(),
		SQS:                     mockSqs,
		Name:                    "test",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		ReceiveQueueName:        "receive",
		SendQueueName:           "send",
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ServiceForwardTestSuite) TestServiceForward_FailOnGetUrl() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{connItemOne()},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	mockSqs.On("GetQueueUrl", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, nil)

	err := Handler(&HandleMessageInput{
		Message:                 sqsMessage(),
		SQS:                     mockSqs,
		Name:                    "test",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		ReceiveQueueName:        "receive",
		SendQueueName:           "send",
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSqs.AssertNotCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func (suite *ServiceForwardTestSuite) TestServiceForward_FailOnDeleteMessage() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)
	mockSqs := new(mocks.SQSAPI)

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{connItemOne()},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	url := "url"
	mockSqs.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: &url,
	}, nil)

	mockSqs.On("DeleteMessage", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(&HandleMessageInput{
		Message:                 sqsMessage(),
		SQS:                     mockSqs,
		Name:                    "test",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		ReceiveQueueName:        "receive",
		SendQueueName:           "send",
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSqs.AssertCalled(suite.T(), "DeleteMessage", mock.Anything)
}

func TestServiceForwardTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceForwardTestSuite))
}
