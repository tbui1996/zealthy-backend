// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	router "github.com/circulohealth/sonar-backend/packages/common/router"
	mock "github.com/stretchr/testify/mock"
)

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// Send provides a mock function with given fields: input
func (_m *Router) Send(input *router.RouterSendInput) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(*router.RouterSendInput) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
