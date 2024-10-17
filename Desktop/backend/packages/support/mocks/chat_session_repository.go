// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	iface "github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	mock "github.com/stretchr/testify/mock"

	model "github.com/circulohealth/sonar-backend/packages/support/pkg/model"

	request "github.com/circulohealth/sonar-backend/packages/support/pkg/request"
)

// ChatSessionRepository is an autogenerated mock type for the ChatSessionRepository type
type ChatSessionRepository struct {
	mock.Mock
}

// AssignPending provides a mock function with given fields: sessionId, internalUserId
func (_m *ChatSessionRepository) AssignPending(sessionId int, internalUserId string) (iface.ChatSession, error) {
	ret := _m.Called(sessionId, internalUserId)

	var r0 iface.ChatSession
	if rf, ok := ret.Get(0).(func(int, string) iface.ChatSession); ok {
		r0 = rf(sessionId, internalUserId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(sessionId, internalUserId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: entity
func (_m *ChatSessionRepository) Create(entity *request.ChatSessionCreateRequest) (iface.ChatSession, error) {
	ret := _m.Called(entity)

	var r0 iface.ChatSession
	if rf, ok := ret.Get(0).(func(*request.ChatSessionCreateRequest) iface.ChatSession); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*request.ChatSessionCreateRequest) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePending provides a mock function with given fields: entity
func (_m *ChatSessionRepository) CreatePending(entity *model.PendingChatSessionCreate) (iface.ChatSession, error) {
	ret := _m.Called(entity)

	var r0 iface.ChatSession
	if rf, ok := ret.Get(0).(func(*model.PendingChatSessionCreate) iface.ChatSession); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(create *model.PendingChatSessionCreate) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntities provides a mock function with given fields: chatType
func (_m *ChatSessionRepository) GetEntities(userID string, chatType string) ([]iface.ChatSession, error) {
	ret := _m.Called(chatType)

	var r0 []iface.ChatSession
	if rf, ok := ret.Get(0).(func(string) []iface.ChatSession); ok {
		r0 = rf(chatType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(chatType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntitiesByExternalID provides a mock function with given fields: userId
func (_m *ChatSessionRepository) GetEntitiesByExternalID(userId string) ([]iface.ChatSession, error) {
	ret := _m.Called(userId)

	var r0 []iface.ChatSession
	if rf, ok := ret.Get(0).(func(string) []iface.ChatSession); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntityWithStatus provides a mock function with given fields: sessionId
func (_m *ChatSessionRepository) GetEntityWithStatus(sessionId string) (iface.ChatSession, *model.ChatStatus, error) {
	ret := _m.Called(sessionId)

	var r0 iface.ChatSession
	if rf, ok := ret.Get(0).(func(string) iface.ChatSession); ok {
		r0 = rf(sessionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(iface.ChatSession)
		}
	}

	var r1 *model.ChatStatus
	if rf, ok := ret.Get(1).(func(string) *model.ChatStatus); ok {
		r1 = rf(sessionId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.ChatStatus)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(sessionId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetEntityWithUsers provides a mock function with given fields: sessionId
func (_m *ChatSessionRepository) GetEntityWithUsers(sessionId string) (iface.ChatSession, error) {
	ret := _m.Called(sessionId)

	var r0 iface.ChatSession
	if rf, ok := ret.Get(0).(func(string) iface.ChatSession); ok {
		r0 = rf(sessionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(iface.ChatSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleReadReceipt provides a mock function with given fields: session, readTime, _a2
func (_m *ChatSessionRepository) HandleReadReceipt(session iface.ChatSession, readTime int64, _a2 request.ReadReceiptRequest) error {
	ret := _m.Called(session, readTime, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(iface.ChatSession, int64, request.ReadReceiptRequest) error); ok {
		r0 = rf(session, readTime, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
