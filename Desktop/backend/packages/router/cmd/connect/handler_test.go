package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
)

type ConnectTestSuite struct {
	suite.Suite
}

func (suite *ConnectTestSuite) TestConnect_Success() {
	mockDB := new(mocks.DynamoDBAPI)
	mockEventPublisher := new(mocks.EventPublisher)

	req := ConnectRequest{
		ConnectionId:   "1",
		UserID:         "test@circulohealth.com",
		Dynamo:         mockDB,
		EventPublisher: mockEventPublisher,
	}

	av, _ := dynamodbattribute.MarshalMap(dynamo.ConnectionItem{
		ConnectionId: req.ConnectionId,
		UserID:       req.UserID,
	})

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarWebsocketConnections),
	}

	mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
	mockEventPublisher.On("PublishConnectionCreatedEvent", mock.Anything, mock.Anything).Return(nil)

	err := Handler(req)

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "PutItem", input)
	mockEventPublisher.AssertCalled(suite.T(), "PublishConnectionCreatedEvent", "test@circulohealth.com", eventconstants.ROUTER_SERVICE)
}

func (suite *ConnectTestSuite) TestConnect_PersistFail() {
	mockDB := new(mocks.DynamoDBAPI)
	mockEventPublisher := new(mocks.EventPublisher)

	req := ConnectRequest{
		ConnectionId:   "1",
		UserID:         "test@circulohealth.com",
		Dynamo:         mockDB,
		EventPublisher: mockEventPublisher,
	}

	av, _ := dynamodbattribute.MarshalMap(dynamo.ConnectionItem{
		ConnectionId: req.ConnectionId,
		UserID:       req.UserID,
	})

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarWebsocketConnections),
	}

	mockDB.On("PutItem", mock.Anything).Return(nil, errors.New("FAKE UNIT TEST ERROR"))
	mockEventPublisher.On("PublishConnectionCreatedEvent", mock.Anything, mock.Anything).Return(nil)

	err := Handler(req)

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "PutItem", input)
}

func (suite *ConnectTestSuite) TestConnect_EventFail() {
	mockDB := new(mocks.DynamoDBAPI)
	mockEventPublisher := new(mocks.EventPublisher)

	req := ConnectRequest{
		ConnectionId:   "1",
		UserID:         "test@circulohealth.com",
		Dynamo:         mockDB,
		EventPublisher: mockEventPublisher,
	}

	mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
	mockEventPublisher.On("PublishConnectionCreatedEvent", mock.Anything, mock.Anything).Return(fmt.Errorf("uh-oh"))

	err := Handler(req)

	suite.NotNil(err)
	mockEventPublisher.AssertCalled(suite.T(), "PublishConnectionCreatedEvent", "test@circulohealth.com", eventconstants.ROUTER_SERVICE)
}

/* Execute Suites */

func TestConnectTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectTestSuite))
}
