package main

import (
	"fmt"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type RecordOfflineMessageHandlerSuite struct {
	suite.Suite
}

func (suite *RecordOfflineMessageHandlerSuite) TestCreatedSuccess() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())

	repo.On("Create", "test-123").Return(true, nil)

	dto, err := Handler("test-123", HandlerDependencies{
		logger,
		repo,
	})

	repo.AssertCalled(suite.T(), "Create", "test-123")
	suite.Nil(err)
	suite.True(dto.Created)
}

func (suite *RecordOfflineMessageHandlerSuite) TestNotCreatedSuccess() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())

	repo.On("Create", "test-123").Return(false, nil)

	dto, err := Handler("test-123", HandlerDependencies{
		logger,
		repo,
	})

	repo.AssertCalled(suite.T(), "Create", "test-123")
	suite.Nil(err)
	suite.False(dto.Created)
}

func (suite *RecordOfflineMessageHandlerSuite) TestFailure() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())

	repo.On("Create", "test-123").Return(false, fmt.Errorf("error"))

	dto, err := Handler("test-123", HandlerDependencies{
		logger,
		repo,
	})

	repo.AssertCalled(suite.T(), "Create", "test-123")
	suite.NotNil(err)
	suite.False(dto.Created)
}

func TestRecordOfflineMessageHandler(t *testing.T) {
	suite.Run(t, new(RecordOfflineMessageHandlerSuite))
}
