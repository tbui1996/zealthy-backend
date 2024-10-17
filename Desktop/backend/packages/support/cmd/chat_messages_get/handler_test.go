package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model/modeltest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
}

func (s *HandlerSuite) Test_HandlesGetMessagesError() {
	sessionID := "test-id"
	mockRepo := new(mocks.ChatMessageRepository)
	mockLogger := zaptest.NewLogger(s.T())

	mockRepo.On("GetMessagesForSession", sessionID).Return(nil, fmt.Errorf("ERROR DUDE!!"))

	_, err := Handler(mockLogger, mockRepo, sessionID)

	mockRepo.AssertCalled(s.T(), "GetMessagesForSession", sessionID)
	s.Equal(err.Error(), "error getting chat messages from dynamo: test-id (ERROR DUDE!!)")
}

func (s *HandlerSuite) Test_ReturnsTheMessagesAsJson() {
	sessionID := "test-id"
	mockRepo := new(mocks.ChatMessageRepository)
	mockLogger := zaptest.NewLogger(s.T())
	messages := new(modeltest.ChatMessageBuilder).
		WithLen(3).
		WithSessionId(sessionID).
		Build()

	mockRepo.On("GetMessagesForSession", sessionID).Return(messages, nil)

	result, _ := Handler(mockLogger, mockRepo, sessionID)
	var parsedResult []model.ChatMessage
	_ = json.Unmarshal([]byte(result), &parsedResult)

	s.Equal(messages, parsedResult)
}

func TestMessagesGet(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
