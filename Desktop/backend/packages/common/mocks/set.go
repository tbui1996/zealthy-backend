// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	arrayiter "github.com/lestrrat-go/iter/arrayiter"

	jwk "github.com/lestrrat-go/jwx/jwk"

	mock "github.com/stretchr/testify/mock"
)

// Set is an autogenerated mock type for the Set type
type Set struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *Set) Add(_a0 jwk.Key) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(jwk.Key) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Clear provides a mock function with given fields:
func (_m *Set) Clear() {
	_m.Called()
}

// Clone provides a mock function with given fields:
func (_m *Set) Clone() (jwk.Set, error) {
	ret := _m.Called()

	var r0 jwk.Set
	if rf, ok := ret.Get(0).(func() jwk.Set); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwk.Set)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *Set) Get(_a0 int) (jwk.Key, bool) {
	ret := _m.Called(_a0)

	var r0 jwk.Key
	if rf, ok := ret.Get(0).(func(int) jwk.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwk.Key)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(int) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Index provides a mock function with given fields: _a0
func (_m *Set) Index(_a0 jwk.Key) int {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(jwk.Key) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Iterate provides a mock function with given fields: _a0
func (_m *Set) Iterate(_a0 context.Context) arrayiter.Iterator {
	ret := _m.Called(_a0)

	var r0 arrayiter.Iterator
	if rf, ok := ret.Get(0).(func(context.Context) arrayiter.Iterator); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(arrayiter.Iterator)
		}
	}

	return r0
}

// Len provides a mock function with given fields:
func (_m *Set) Len() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// LookupKeyID provides a mock function with given fields: _a0
func (_m *Set) LookupKeyID(_a0 string) (jwk.Key, bool) {
	ret := _m.Called(_a0)

	var r0 jwk.Key
	if rf, ok := ret.Get(0).(func(string) jwk.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwk.Key)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: _a0
func (_m *Set) Remove(_a0 jwk.Key) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(jwk.Key) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
