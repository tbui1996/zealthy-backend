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

type DisconnectTestSuite struct {
	suite.Suite
	input *dynamodb.DeleteItemInput
}

func (suite *DisconnectTestSuite) SetupSuite() {
	suite.input = &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String("1"),
			},
			"UserID": {
				S: aws.String("circulo"),
			},
		},
		TableName: aws.String(dynamo.SonarWebsocketConnections),
	}
}

func (suite *DisconnectTestSuite) TestDisconnect_Success() {
	mockDB := new(mocks.DynamoDBAPI)

	mockDB.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := Handler(DisconnectRequest{
		ConnectionId: "1",
		UserId:       "circulo",
		Dynamo:       mockDB,
	})

	suite.Nil(err)
	mockDB.AssertCalled(suite.T(), "DeleteItem", suite.input)
}

func (suite *DisconnectTestSuite) TestDisconnect_Fail() {
	mockDB := new(mocks.DynamoDBAPI)

	mockDB.On("DeleteItem", mock.Anything).Return(nil, errors.New("FAKE UNIT TEST ERROR"))

	err := Handler(DisconnectRequest{
		ConnectionId: "1",
		UserId:       "circulo",
		Dynamo:       mockDB,
	})

	suite.NotNil(err)
	mockDB.AssertCalled(suite.T(), "DeleteItem", suite.input)
}

/* Execute Suites */

func TestDisconnectTestSuite(t *testing.T) {
	suite.Run(t, new(DisconnectTestSuite))
}
