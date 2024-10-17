package iface

import (
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
)

type ChatSession interface {
	ID() string
	InternalUserID() string
	Patient() response.PatientResponse
	UserID() string
	IsPending() bool
	LastMessageTimestamp() int64
	UserLastRead() int64
	InternalUserLastRead() int64
	ToResponseDTO() response.ChatSessionResponseDTO
	GetMessages() ([]model.ChatMessage, error)
	AppendRequestMessage(message request.Chat) (*model.ChatMessage, error)
	Type() model.ChatType
}
