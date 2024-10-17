package response

import (
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/stretchr/testify/suite"
)

type ChatSessionJoinResponseTestSuite struct {
	suite.Suite
	openResponse    ChatSessionJoinResponse
	openResponse1   ChatSessionJoinResponse
	closeResponse   ChatSessionJoinResponse
	closeResponse1  ChatSessionJoinResponse
	pendingResponse ChatSessionJoinResponse
}

func (suite *ChatSessionJoinResponseTestSuite) SetupTest() {
	lastSent := 0
	lastRead := 0
	lastMessage := "Test"

	suite.openResponse = ChatSessionJoinResponse{
		ID:          1,
		Created:     time.Now(),
		Status:      model.OPEN,
		UserId:      "Olive",
		LastRead:    &lastRead,
		Name:        "TOPIC",
		Value:       "Test User",
		LastMessage: &lastMessage,
		LastSent:    &lastSent,
	}

	tempResp := suite.openResponse
	tempResp.UserId = "Okta"
	suite.openResponse1 = tempResp

	tempResp = suite.openResponse
	tempResp.ID = 2
	tempResp.Status = model.CLOSED
	suite.closeResponse = tempResp
	tempResp.UserId = "Okta"
	suite.closeResponse1 = tempResp

	tempResp = suite.openResponse
	tempResp.ID = 3
	tempResp.Status = model.PENDING
	suite.pendingResponse = tempResp
}

func (suite *ChatSessionJoinResponseTestSuite) TestNormalize_SingleResultNoPending() {
	responses := ChatSessionJoinResponses{suite.openResponse, suite.openResponse1}

	dto, pending := responses.NormalizeSingleResultToDTO()

	suite.NotNil(dto)
	suite.NotEmpty(dto.Topic)
	suite.True(dto.ChatOpen)
	suite.False(pending)
}

func (suite *ChatSessionJoinResponseTestSuite) TestNormalize_MultipleResultNoPending() {
	responses := ChatSessionJoinResponses{suite.openResponse, suite.openResponse1, suite.closeResponse, suite.closeResponse1}

	dto := responses.NormalizeMultipleResultsToDTO()

	suite.Len(dto, 2)
}

func (suite *ChatSessionJoinResponseTestSuite) TestNormalize_SingleResultWithPending() {
	responses := ChatSessionJoinResponses{suite.pendingResponse}

	dto, pending := responses.NormalizeSingleResultToDTO()

	suite.NotNil(dto)
	suite.True(dto.ChatOpen)
	suite.NotEmpty(dto.Topic)
	suite.True(pending)
}

func (suite *ChatSessionJoinResponseTestSuite) TestNormalize_MultipleResultWithPending() {
	responses := ChatSessionJoinResponses{suite.pendingResponse, suite.openResponse, suite.openResponse1, suite.closeResponse, suite.closeResponse1}

	dto := responses.NormalizeMultipleResultsToDTO()

	suite.Len(dto, 3)
}

func (suite *ChatSessionJoinResponseTestSuite) TestNormalize_MultipleResultKeepsOrder() {
	responses := ChatSessionJoinResponses{suite.pendingResponse, suite.openResponse, suite.openResponse1, suite.closeResponse, suite.closeResponse1}
	dto := responses.NormalizeMultipleResultsToDTO()
	suite.Len(dto, 3)
	suite.Equal("3", dto[0].ChatSessionDTO.ID)
	suite.Equal("1", dto[1].ChatSessionDTO.ID)
	suite.Equal("2", dto[2].ChatSessionDTO.ID)
}

func TestChatSessionJoinResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ChatSessionJoinResponseTestSuite))
}
