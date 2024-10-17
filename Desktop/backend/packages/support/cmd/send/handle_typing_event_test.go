package main

import (
	"encoding/json"
	"errors"
	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandleTypingTestSuite struct {
	suite.Suite
	config requestConfig.APIGatewayWebsocketProxyRequest
	req    request.TypingActionRequest
	repo   *mocks.ChatSessionRepository
	client *router.MockRouter
	sess   *mocks.ChatSession
}

func (suite *HandleTypingTestSuite) SetupTest() {
	suite.config = requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(suite.T()),
		Event: events.APIGatewayWebsocketProxyRequest{
			RequestContext: events.APIGatewayWebsocketProxyRequestContext{
				Authorizer: map[string]interface{}{}},
		},
	}

	suite.repo = new(mocks.ChatSessionRepository)

	suite.req = request.TypingActionRequest{
		SessionID: "1",
		UserID:    "1",
		Action:    "start",
	}

	suite.client = new(router.MockRouter)

	suite.sess = new(mocks.ChatSession)
}

func (suite *HandleTypingTestSuite) TestHandleTyping_SendsMessage() {
	message, _ := json.Marshal(suite.req)

	suite.repo.On("GetEntityWithUsers", "1").Return(suite.sess, nil)

	suite.sess.On("UserID").Return("1")

	suite.client.On("Send", mock.Anything).Return(nil)

	mockClient := &router.Session{
		Router: suite.client,
	}

	resErr := HandleTyping(&suite.config, string(message), suite.repo, mockClient)

	suite.Nil(resErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", "1")
	suite.sess.AssertCalled(suite.T(), "UserID")
	suite.client.AssertCalled(suite.T(), "Send", mock.Anything)
}

func (suite *HandleTypingTestSuite) TestHandleTyping_FailsGetSession() {
	message, _ := json.Marshal(suite.req)

	suite.repo.On("GetEntityWithUsers", "1").Return(nil, errors.New("FAKE ERROR"))

	mockClient := &router.Session{
		Router: suite.client,
	}

	resErr := HandleTyping(&suite.config, string(message), suite.repo, mockClient)

	suite.Error(resErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", "1")
	suite.sess.AssertNotCalled(suite.T(), "UserID")
	suite.client.AssertNotCalled(suite.T(), "Send", mock.Anything)
}

func (suite *HandleTypingTestSuite) TestHandleTyping_FailsSend() {
	message, _ := json.Marshal(suite.req)

	suite.repo.On("GetEntityWithUsers", "1").Return(suite.sess, nil)

	suite.sess.On("UserID").Return("1")

	suite.client.On("Send", mock.Anything).Return(errors.New("FAKE ERROR"))

	mockClient := &router.Session{
		Router: suite.client,
	}

	resErr := HandleTyping(&suite.config, string(message), suite.repo, mockClient)

	suite.Error(resErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", "1")
	suite.sess.AssertCalled(suite.T(), "UserID")
	suite.client.AssertCalled(suite.T(), "Send", mock.Anything)
}

func TestHandleTypingTestSuite(t *testing.T) {
	suite.Run(t, new(HandleTypingTestSuite))
}
