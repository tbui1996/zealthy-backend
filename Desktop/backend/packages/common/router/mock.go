package router

import (
	"github.com/stretchr/testify/mock"
)

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) Send(input *RouterSendInput) error {
	args := m.Called(input)

	if args.Get(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}
