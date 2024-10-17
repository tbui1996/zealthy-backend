package main

import (
	"encoding/json"
	"fmt"
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
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)

	mockRepo.On("GetEntities", "CIRCULATOR").Return(nil, fmt.Errorf("REKT"))

	_, err := Handler("1", "internals_program_manager", mockLogger, mockRepo)

	s.Equal("unable to get connection items REKT", err.Error())
}

func (s *HandlerSuite) TestHandlesInvalidTypes() {
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)

	mockRepo.On("GetEntities", "CIRCULATOR").Return(nil, fmt.Errorf("REKT"))

	_, err := Handler("1", "internals_non_role", mockLogger, mockRepo)

	s.Error(err)
}

func (s *HandlerSuite) TestConvertsToResponseDtos() {
	mockLogger := zaptest.NewLogger(s.T())
	mockRepo := new(mocks.ChatSessionRepository)
	mockSession := new(mocks.ChatSession)
	mockRepoReturn := []iface.ChatSession{
		mockSession,
		mockSession,
		mockSession,
	}
	fakeDto := response.ChatSessionResponseDTO{ID: "SOMETHING RANDOM"}
	expectedResult := []response.ChatSessionResponseDTO{
		fakeDto,
		fakeDto,
		fakeDto,
	}

	mockSession.On("ToResponseDTO").Return(fakeDto)
	mockRepo.On("GetEntities", "CIRCULATOR").Return(mockRepoReturn, nil)

	result, _ := Handler("1", "internals_program_manager", mockLogger, mockRepo)
	var parsedResult []response.ChatSessionResponseDTO
	_ = json.Unmarshal(result, &parsedResult)

	s.Equal(expectedResult, parsedResult)
}

func TestChatSessionsGet(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
