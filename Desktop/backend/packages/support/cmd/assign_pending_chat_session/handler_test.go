package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type MainSuite struct {
	suite.Suite
}

func (s *MainSuite) Test_HandlesRepoError() {
	logger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)
	request := request.AssignPendingChatSessionRequestInternal{SessionID: "123", InternalUserID: "321"}

	mockRepo.On("AssignPending", 123, "321").Return(nil, fmt.Errorf("error"))

	_, err := Handler(logger, mockRepo, request)
	s.Equal(err.StatusCode, http.StatusInternalServerError)
	s.Equal(err.Error(), "error storing item: (error)")
}

func (s *MainSuite) Test_HandlesMessagesError() {
	logger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)
	mockSession := new(mocks.ChatSession)
	request := request.AssignPendingChatSessionRequestInternal{SessionID: "123", InternalUserID: "321"}

	mockRepo.On("AssignPending", 123, "321").Return(mockSession, nil)
	mockSession.On("GetMessages").Return(nil, fmt.Errorf("oh noes!"))

	_, err := Handler(logger, mockRepo, request)

	s.Equal(err.StatusCode, http.StatusInternalServerError)
	s.Equal(err.Error(), "error getting pending chat messages: 123 (oh noes!)")
}

func (s *MainSuite) Test_ReturnsCorrectJson() {
	logger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)
	mockSession := new(mocks.ChatSession)
	messages := []model.ChatMessage{{Message: "Hi"}, {Message: "Hello!"}, {Message: "How was your weekend?"}}
	request := request.AssignPendingChatSessionRequestInternal{SessionID: "123", InternalUserID: "321"}

	expected, _ := json.Marshal(messages)

	mockRepo.On("AssignPending", 123, "321").Return(mockSession, nil)
	mockSession.On("GetMessages").Return(messages, nil)

	result, _ := Handler(logger, mockRepo, request)

	s.Equal(expected, result)
}

func TestAssignPendingChatSession(t *testing.T) {
	suite.Run(t, new(MainSuite))
}
