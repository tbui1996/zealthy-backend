package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnconfirmedConnectTestSuite struct {
	suite.Suite
}

func (suite *UnconfirmedConnectTestSuite) TestUnconfirmedConnect_Success() {
	mockDB := new(mocks.DynamoDBAPI)

	req := UnconfirmedConnectRequest{
		ConnectionId: "1",
		Email:        "test@circulohealth.com",
		Dynamo:       mockDB,
	}

	av, _ := dynamodbattribute.MarshalMap(dynamo.UnconfirmedConnectionItem{
		ConnectionId: req.ConnectionId,
		Email:        req.Email,
	})

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	mockDB.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	err := Handler(UnconfirmedConnectRequest{
		ConnectionId: "1",
		Email:        "test@circulohealth.com",
		Dynamo:       mockDB,
	})

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "PutItem", input)
}

func (suite *UnconfirmedConnectTestSuite) TestUnconfirmedConnect_Fail() {
	mockDB := new(mocks.DynamoDBAPI)

	req := UnconfirmedConnectRequest{
		ConnectionId: "1",
		Email:        "test@circulohealth.com",
		Dynamo:       mockDB,
	}

	av, _ := dynamodbattribute.MarshalMap(dynamo.UnconfirmedConnectionItem{
		ConnectionId: req.ConnectionId,
		Email:        req.Email,
	})

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	mockDB.On("PutItem", mock.Anything).Return(nil, errors.New("FAKE UNIT TEST ERROR"))

	err := Handler(req)

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "PutItem", input)
}

/* Execute Suites */

func TestUnconfirmedConnectTestSuite(t *testing.T) {
	suite.Run(t, new(UnconfirmedConnectTestSuite))
}
