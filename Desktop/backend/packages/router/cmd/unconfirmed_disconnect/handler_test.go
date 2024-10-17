package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnconfirmedDisconnectTestSuite struct {
	suite.Suite
}

func (suite *UnconfirmedDisconnectTestSuite) TestUnconfirmedDisconnect_Success() {
	mockDB := new(mocks.DynamoDBAPI)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String("1"),
			},
			"Email": {
				S: aws.String("test@circulohealth.com"),
			},
		},
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	mockDB.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := Handler(UnconfirmedDisconnectRequest{
		ConnectionId: "1",
		Email:        "test@circulohealth.com",
		Dynamo:       mockDB,
	})

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "DeleteItem", input)
}

func (suite *UnconfirmedDisconnectTestSuite) TestUnconfirmedDisconnect_Fail() {
	mockDB := new(mocks.DynamoDBAPI)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String("1"),
			},
			"Email": {
				S: aws.String("test@circulohealth.com"),
			},
		},
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	mockDB.On("DeleteItem", mock.Anything).Return(nil, errors.New("FAKE UNIT TEST ERROR"))

	err := Handler(UnconfirmedDisconnectRequest{
		ConnectionId: "1",
		Email:        "test@circulohealth.com",
		Dynamo:       mockDB,
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "DeleteItem", input)
}

/* Execute Suites */

func TestUnconfirmedDisconnectTestSuite(t *testing.T) {
	suite.Run(t, new(UnconfirmedDisconnectTestSuite))
}
