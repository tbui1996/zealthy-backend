package session

import (
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/mock"
)

type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) Create(message request.Chat) (*model.ChatMessage, error) {
	args := m.Called(message)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) GetMessagesForSession(id string) ([]model.ChatMessage, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.ChatMessage), args.Error(1)
}
