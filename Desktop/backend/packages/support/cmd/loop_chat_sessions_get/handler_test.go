package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
}

func (s *HandlerSuite) TestHandlesRepoError() {
	logger := zaptest.NewLogger(s.T())
	repo := new(mocks.ChatSessionRepository)
	userID := "123"

	repo.On("GetEntitiesByExternalID", userID).Return(nil, fmt.Errorf("uh-oh"))

	_, err := Handler(logger, repo, userID)
	s.Equal(http.StatusInternalServerError, err.StatusCode)
	s.Equal("failed to get chat sessions for user. uh-oh", err.Error())
}

func (s *HandlerSuite) TestReturnsDTOS() {
	logger := zaptest.NewLogger(s.T())
	repo := new(mocks.ChatSessionRepository)
	session := new(mocks.ChatSession)
	sessions := []iface.ChatSession{
		session,
	}
	userID := "123"
	testDTO := response.ChatSessionResponseDTO{ID: "TEST-ID", UserID: userID}
	expected := []response.ChatSessionResponseDTO{testDTO}

	repo.On("GetEntitiesByExternalID", userID).Return(sessions, nil)
	session.On("ToResponseDTO").Return(testDTO)

	result, _ := Handler(logger, repo, userID)

	var parsedResult []response.ChatSessionResponseDTO
	_ = json.Unmarshal(result, &parsedResult)

	s.Equal(expected, parsedResult)
}

func TestLoopChatSessionGet(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
