package main

import (
	"errors"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"testing"
)

type GetOrganizationsSuite struct {
	suite.Suite
	input                          GetOrganizationsRequest
	externalUserOrganizationMapper *iface.MockExternalUserOrganization
	registry                       *mapper.MockRegistryAPI
	validOrgs                      []*model.ExternalUserOrganization
	orgsWithNil                    []*model.ExternalUserOrganization
}

func (suite *GetOrganizationsSuite) SetupTest() {
	suite.registry = new(mapper.MockRegistryAPI)
	suite.externalUserOrganizationMapper = new(iface.MockExternalUserOrganization)

	suite.input = GetOrganizationsRequest{
		Registry: suite.registry,
		Logger:   zaptest.NewLogger(suite.T()),
	}

	suite.validOrgs = []*model.ExternalUserOrganization{
		{ID: 1},
		{ID: 2},
	}

	suite.orgsWithNil = []*model.ExternalUserOrganization{
		{ID: 1},
		nil,
		{ID: 3},
	}
}

func (suite *GetOrganizationsSuite) TestGetOrganizations_Success() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("FindAll").Return(suite.validOrgs, nil)

	actualOut, actualErr := handler(suite.input)

	suite.NotNil(actualOut)
	suite.Equal(2, len(actualOut))
	suite.NoError(actualErr)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "FindAll")
}

func (suite *GetOrganizationsSuite) TestGetOrganizations_SuccessWithNil() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("FindAll").Return(suite.orgsWithNil, nil)

	actualOut, actualErr := handler(suite.input)

	suite.NotNil(actualOut)
	suite.Equal(2, len(actualOut))
	suite.NoError(actualErr)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "FindAll")
}

func (suite *GetOrganizationsSuite) TestGetOrganizations_Fail() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("FindAll").Return(nil, errors.New("FAKE TEST ERROR IGNORE"))

	actualOut, actualErr := handler(suite.input)

	suite.Nil(actualOut)
	suite.Error(actualErr)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "FindAll")
}

func (suite *GetOrganizationsSuite) TestHandleError() {
	errMsg := errors.New("FAKE ERROR, IGNORE IT")
	resp, _ := errorResponse(errMsg, "THIS IS A", zaptest.NewLogger(suite.T()))

	suite.Equal("THIS IS A FAKE ERROR, IGNORE IT", resp.Body)
}

func TestExternalSignUp(t *testing.T) {
	suite.Run(t, new(GetOrganizationsSuite))
}
