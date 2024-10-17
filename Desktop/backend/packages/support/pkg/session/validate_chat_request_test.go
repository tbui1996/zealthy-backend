package session

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
)

/* Mocks */
type MockDatabase struct{}

func (m *MockDatabase) Create(item interface{}) error {
	return nil
}

func (m *MockDatabase) Update(expression string, item interface{}, key interface{}) error {
	return nil
}

type ValidateChatRequestTestSuite struct {
	suite.Suite
}

func (suite *ValidateChatRequestTestSuite) TestValidRequestDoesNotReturnError() {
	err := validateChatRequest(request.Chat{
		Session: "1",
		Sender:  "1",
		Message: "test",
		File:    "",
	})

	suite.Nil(err)
}

func (suite *ValidateChatRequestTestSuite) TestMissingSessionReturnsError() {
	err := validateChatRequest(request.Chat{
		Sender:  "1",
		Message: "test",
		File:    "",
	})

	suite.True(suite.Error(err))
}

func (suite *ValidateChatRequestTestSuite) TestMissingSenderReturnsError() {
	err := validateChatRequest(request.Chat{
		Session: "1",
		Message: "test",
		File:    "",
	})

	suite.True(suite.Error(err))
}

func (suite *ValidateChatRequestTestSuite) TestMissingMessageReturnsError() {
	err := validateChatRequest(request.Chat{
		Session: "1",
		Sender:  "1",
		File:    "",
	})

	suite.True(suite.Error(err))
}

/* Execute Suites */
func TestValidateChatRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateChatRequestTestSuite))
}
