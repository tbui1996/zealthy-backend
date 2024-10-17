package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/forward"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BroadcastTestSuite struct {
	suite.Suite
}

func (suite *BroadcastTestSuite) TestBroadcast_Success() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	item := map[string]*dynamodb.AttributeValue{
		"ConnectionId": {
			S: aws.String("1"),
		},
		"UserID": {
			S: aws.String("1"),
		},
	}

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{item},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	err := Handler(forward.ForwarderBroadcastDTO{
		Sender:                  "1",
		Recipients:              []string{"1"},
		Message:                 "message",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *BroadcastTestSuite) TestBroadcastScan_Fail() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("Scan", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(forward.ForwarderBroadcastDTO{
		Sender:                  "1",
		Recipients:              []string{"1"},
		Message:                 "message",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *BroadcastTestSuite) TestBroadcastPostConnection_Fail() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	item := map[string]*dynamodb.AttributeValue{
		"ConnectionId": {
			S: aws.String("1"),
		},
		"UserID": {
			S: aws.String("1"),
		},
	}

	mockDB.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{item},
	}, nil)

	mockApi.On("PostToConnection", mock.Anything).Return(
		nil,
		errors.New("FAKE TEST ERROR, IGNORE"),
	)

	err := Handler(forward.ForwarderBroadcastDTO{
		Sender:                  "1",
		Recipients:              []string{"1"},
		Message:                 "message",
		DynamoDB:                mockDB,
		ApiGatewayManagementApi: mockApi,
		Logger:                  zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "Scan", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func TestBroadcastTestSuite(t *testing.T) {
	suite.Run(t, new(BroadcastTestSuite))
}
