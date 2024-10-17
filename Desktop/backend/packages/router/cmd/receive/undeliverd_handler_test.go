package main

import (
	"errors"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UndeliveredHandlerTestSuite struct {
	suite.Suite
}

func (suite *UndeliveredHandlerTestSuite) TestUndeliveredHandler_Success() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(getQueryOutput(), nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	mockDB.On("DeleteItem", mock.Anything).Return(nil, nil)

	err := UndeliveredHandler(mockDB, mockApi, "1")

	suite.Nil(err)

	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertCalled(suite.T(), "DeleteItem", mock.Anything)
}

func (suite *UndeliveredHandlerTestSuite) TestUndeliveredHandler_FailOnGetItem() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("GetItem", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := UndeliveredHandler(mockDB, mockApi, "1")

	suite.NotNil(err)

	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
}

func (suite *UndeliveredHandlerTestSuite) TestUndeliveredHandler_FailOnQuery() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := UndeliveredHandler(mockDB, mockApi, "1")

	suite.NotNil(err)

	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
}

func (suite *UndeliveredHandlerTestSuite) TestUndeliveredHandler_FailOnPost() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(getQueryOutput(), nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := UndeliveredHandler(mockDB, mockApi, "1")

	suite.Nil(err)

	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertNotCalled(suite.T(), "DeleteItem", mock.Anything)
}

func (suite *UndeliveredHandlerTestSuite) TestUndeliveredHandler_FailOnDelete() {
	mockDB := new(mocks.DynamoDBAPI)
	mockApi := new(mocks.ApiGatewayManagementApiAPI)

	mockDB.On("GetItem", mock.Anything).Return(getItemOutput(), nil)

	mockDB.On("Query", mock.Anything).Return(getQueryOutput(), nil)

	mockApi.On("PostToConnection", mock.Anything).Return(nil, nil)

	mockDB.On("DeleteItem", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := UndeliveredHandler(mockDB, mockApi, "1")

	suite.Nil(err)

	mockDB.AssertCalled(suite.T(), "GetItem", mock.Anything)
	mockDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockApi.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	mockDB.AssertCalled(suite.T(), "DeleteItem", mock.Anything)
}

func TestUndeliveredHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UndeliveredHandlerTestSuite))
}
