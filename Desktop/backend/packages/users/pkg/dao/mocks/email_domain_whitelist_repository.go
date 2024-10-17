// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	exception "github.com/circulohealth/sonar-backend/packages/common/exception"

	mock "github.com/stretchr/testify/mock"

	model "github.com/circulohealth/sonar-backend/packages/users/pkg/model"
)

// EmailDomainWhitelistRepository is an autogenerated mock type for the EmailDomainWhitelistRepository type
type EmailDomainWhitelistRepository struct {
	mock.Mock
}

// GetWhitelistDomain provides a mock function with given fields: domain
func (_m *EmailDomainWhitelistRepository) GetWhitelistDomain(domain string) (*model.EmailDomainWhitelist, *exception.SonarError) {
	ret := _m.Called(domain)

	var r0 *model.EmailDomainWhitelist
	if rf, ok := ret.Get(0).(func(string) *model.EmailDomainWhitelist); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.EmailDomainWhitelist)
		}
	}

	var r1 *exception.SonarError
	if rf, ok := ret.Get(1).(func(string) *exception.SonarError); ok {
		r1 = rf(domain)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*exception.SonarError)
		}
	}

	return r0, r1
}
