package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sesv2"
	common_mocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/constants"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/input"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type TriggerOfflineEmailHandlerSuite struct {
	suite.Suite
}

const TEMPLATE_TEST string = "template-test"
const TEST string = "test"

func (suite *TriggerOfflineEmailHandlerSuite) TestNotUpdatedError() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())
	input := input.UserInfo{
		ID: "TEST-123",
	}
	sesClient := new(common_mocks.SESV2API)
	template := TEMPLATE_TEST
	domain := TEST
	configSet := TEST
	name := []byte{}

	repo.On("UpdateStatus", input.ID, constants.SENT_NOTIFICATION).Return(false, fmt.Errorf("error"))

	err := Handler(input, HandlerDependencies{
		logger,
		repo,
		sesClient,
		template,
		domain,
		configSet,
		name,
	})

	repo.AssertCalled(suite.T(), "UpdateStatus", input.ID, constants.SENT_NOTIFICATION)
	suite.NotNil(err)
}

func (suite *TriggerOfflineEmailHandlerSuite) TestNotUpdatedSuccess() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())
	input := input.UserInfo{
		ID: "TEST-123",
	}
	sesClient := new(common_mocks.SESV2API)
	template := TEMPLATE_TEST
	domain := TEST
	configSet := TEST
	name := []byte{}

	repo.On("UpdateStatus", input.ID, constants.SENT_NOTIFICATION).Return(false, nil)

	err := Handler(input, HandlerDependencies{
		logger,
		repo,
		sesClient,
		template,
		domain,
		configSet,
		name,
	})

	repo.AssertCalled(suite.T(), "UpdateStatus", input.ID, constants.SENT_NOTIFICATION)
	sesClient.AssertNotCalled(suite.T(), "SendEmail", mock.Anything)
	suite.Nil(err)
}

func (suite *TriggerOfflineEmailHandlerSuite) TestUpdatedSuccess() {
	repo := new(mocks.OfflineMessageNotificationRepo)
	logger := zaptest.NewLogger(suite.T())
	input := input.UserInfo{
		ID: "TEST-123",
	}
	sesClient := new(common_mocks.SESV2API)
	template := TEMPLATE_TEST
	domain := TEST
	configSet := TEST
	name := []byte{}

	repo.On("UpdateStatus", input.ID, constants.SENT_NOTIFICATION).Return(true, nil)
	sesClient.On("SendEmail", mock.Anything).Return(&sesv2.SendEmailOutput{}, nil)

	err := Handler(input, HandlerDependencies{
		logger,
		repo,
		sesClient,
		template,
		domain,
		configSet,
		name,
	})

	repo.AssertCalled(suite.T(), "UpdateStatus", input.ID, constants.SENT_NOTIFICATION)
	sesClient.AssertCalled(suite.T(), "SendEmail", mock.Anything)
	suite.Nil(err)
}

func TestTriggerOfflineHandler(t *testing.T) {
	suite.Run(t, new(TriggerOfflineEmailHandlerSuite))
}
