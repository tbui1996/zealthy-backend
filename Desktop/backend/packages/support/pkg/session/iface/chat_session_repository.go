package iface

import (
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
)

type ChatSessionRepository interface {
	GetEntityWithUsers(sessionId string) (ChatSession, error)
	AssignPending(sessionId int, internalUserId string) (ChatSession, error)
	Create(entity *request.ChatSessionCreateRequest) (ChatSession, error)
	CreatePending(entity *model.PendingChatSessionCreate) (ChatSession, error)
	GetEntities(userID string, chatType string) ([]ChatSession, error)
	GetEntitiesByExternalID(userId string) ([]ChatSession, error)
	GetEntityWithStatus(sessionId string) (ChatSession, *model.ChatStatus, error)
	HandleReadReceipt(session ChatSession, readTime int64, request request.ReadReceiptRequest) error
}
