package modeltest

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChatMessageBuilderSuite struct {
	suite.Suite
}

func (s *ChatMessageBuilderSuite) TestWithLen() {
	messages := new(ChatMessageBuilder).
		WithLen(3).
		Build()

	s.Len(messages, 3)
}

func (s *ChatMessageBuilderSuite) TestBuildOneWithSessionId() {
	message := new(ChatMessageBuilder).
		WithSessionId("123").
		BuildOne()
	s.Equal("123", message.SessionID)
}

func (s *ChatMessageBuilderSuite) TestBuilWithSessionId() {
	message := new(ChatMessageBuilder).
		WithLen(1).
		WithSessionId("123").
		Build()[0]
	s.Equal("123", message.SessionID)
}

func (s *ChatMessageBuilderSuite) TestBuildWithoutLen() {
	message := new(ChatMessageBuilder).Build()

	s.Nil(message)
}

func (s *ChatMessageBuilderSuite) TestWithZero() {
	messages := new(ChatMessageBuilder).
		WithLen(0).
		Build()

	s.Len(messages, 0)
}

func TestChatMessageBuilder(t *testing.T) {
	suite.Run(t, new(ChatMessageBuilderSuite))
}
