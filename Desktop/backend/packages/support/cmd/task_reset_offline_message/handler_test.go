package main

import (
	"fmt"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ResetOfflineMessageHandlerSuite struct {
	suite.Suite
}

func (suite *ResetOfflineMessageHandlerSuite) TestCreatedSuccess() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())

	repo.On("Remove", "test-123").Return(nil)

	err := Handler("test-123", HandlerDependencies{
		logger,
		repo,
	})

	repo.AssertCalled(suite.T(), "Remove", "test-123")
	suite.Nil(err)
}

func (suite *ResetOfflineMessageHandlerSuite) TestCreatedError() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())

	repo.On("Remove", "test-123").Return(fmt.Errorf("error"))

	err := Handler("test-123", HandlerDependencies{
		logger,
		repo,
	})

	repo.AssertCalled(suite.T(), "Remove", "test-123")
	suite.NotNil(err)
}

func TestResetOfflineMessageHandler(t *testing.T) {
	suite.Run(t, new(ResetOfflineMessageHandlerSuite))
}
