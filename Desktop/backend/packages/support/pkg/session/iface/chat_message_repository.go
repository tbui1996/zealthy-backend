package iface

import (
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
)

type ChatMessageRepository interface {
	Create(message request.Chat) (*model.ChatMessage, error)
	GetMessagesForSession(id string) ([]model.ChatMessage, error)
}
