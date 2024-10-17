package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
}

func (s *HandlerSuite) Test_HandlesCreateError() {
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)

	mockRepo.On("Create", mock.Anything).Return(nil, fmt.Errorf("not created uh-oh"))

	_, err := Handler(mockLogger, mockRepo, request.ChatSessionRequestInternal{})

	s.Equal(err.StatusCode, http.StatusInternalServerError)
	s.Contains(err.Error(), "error storing item:")
	s.Contains(err.Error(), "(not created uh-oh)")
}

func (s *HandlerSuite) Test_HandlesNonExistingSession() {
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)

	mockRepo.On("Create", mock.Anything).Return(nil, nil)

	_, err := Handler(mockLogger, mockRepo, request.ChatSessionRequestInternal{})
	s.Equal(err.StatusCode, http.StatusInternalServerError)
	s.Equal(err.Error(), "expected session to exist after creating")
}

func (s *HandlerSuite) Test_ReturnsTheSessionResponse() {
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)
	mockSession := new(mocks.ChatSession)
	input := request.ChatSessionRequestInternal{InternalUserID: "INTERNAL_USER", UserID: "USER"}
	requestMatcher := func(createRequest *request.ChatSessionCreateRequest) bool {
		return *(createRequest.InternalUserID) == input.InternalUserID && createRequest.UserID == input.UserID && createRequest.ChatOpen
	}

	mockRepo.On("Create", mock.MatchedBy(requestMatcher)).Return(mockSession, nil)
	mockSession.On("ID").Return("test-session-id")

	result, _ := Handler(mockLogger, mockRepo, input)

	var response response.ChatSessionResponse
	_ = json.Unmarshal(result, &response)

	s.Equal(response.ID, "test-session-id")
}
func TestChatSessionCreate(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
