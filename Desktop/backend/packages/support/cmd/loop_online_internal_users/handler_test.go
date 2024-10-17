package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type LoopOnlineUsersTestSuite struct {
	suite.Suite
}

func (s *LoopOnlineUsersTestSuite) Test__SuccessDoesNotReturnDuplicates() {
	mockDb := new(mocks.DynamoDBAPI)
	mockLogger := zaptest.NewLogger(s.T())
	existingConnectionItems := []model.ConnectionItem{
		{ConnectionId: "123", UserID: "123"},
		{ConnectionId: "1234", UserID: "123"},
		{ConnectionId: "122", UserID: "12"},
	}

	expected := []response.OnlineUserResponse{
		{UserID: "123"},
		{UserID: "12"},
	}

	var existingConnectionItemMaps = make([]map[string]*dynamodb.AttributeValue, 0)

	for _, item := range existingConnectionItems {
		mapItem, _ := dynamodbattribute.MarshalMap(item)
		existingConnectionItemMaps = append(existingConnectionItemMaps, mapItem)
	}

	mockDb.On("Scan", mock.Anything).Return(&dynamodb.ScanOutput{
		Items: existingConnectionItemMaps,
	}, nil)

	deps := HandlerDeps{
		mockDb,
		mockLogger,
	}

	result, err := Handler(deps)

	s.Equal(&expected, result)
	s.Nil(err)
}

func (s *LoopOnlineUsersTestSuite) Test__Error() {
	mockDb := new(mocks.DynamoDBAPI)
	mockLogger := zaptest.NewLogger(s.T())

	mockDb.On("Scan", mock.Anything).Return(nil, fmt.Errorf("fake error"))

	deps := HandlerDeps{
		mockDb,
		mockLogger,
	}

	_, err := Handler(deps)

	s.NotNil(err)
	s.Equal("fake error", err.Error())
}

func TestOnlineUsersHandlerTests(t *testing.T) {
	suite.Run(t, new(LoopOnlineUsersTestSuite))
}
