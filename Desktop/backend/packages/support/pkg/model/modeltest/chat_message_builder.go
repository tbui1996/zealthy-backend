package modeltest

import (
	"fmt"
	"strconv"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
)

type ChatMessageBuilder struct {
	messages  []model.ChatMessage
	sessionId string
}

func (builder *ChatMessageBuilder) WithLen(length int) *ChatMessageBuilder {
	builder.messages = make([]model.ChatMessage, length)
	return builder
}

func (builder *ChatMessageBuilder) WithSessionId(sessionID string) *ChatMessageBuilder {
	builder.sessionId = sessionID
	return builder
}

func (builder *ChatMessageBuilder) Build() []model.ChatMessage {
	if builder.messages == nil {
		return nil
	}

	for i := 0; i < len(builder.messages); i++ {
		builder.messages[i] = model.ChatMessage{SessionID: builder.sessionId, Message: fmt.Sprintf("message %s", strconv.Itoa(i+1))}
	}

	return builder.messages
}

func (builder *ChatMessageBuilder) BuildOne() *model.ChatMessage {
	return &model.ChatMessage{SessionID: builder.sessionId, Message: "message 1"}
}
