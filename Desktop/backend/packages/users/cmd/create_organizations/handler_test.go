package main

import (
	"errors"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CreateOrganizationsSuite struct {
	suite.Suite
	input                          CreateOrganizationsRequest
	externalUserOrganizationMapper *iface.MockExternalUserOrganization
	registry                       *mapper.MockRegistryAPI
	validOrgs                      []*model.ExternalUserOrganization
	duplicateOrgs                  []*model.ExternalUserOrganization
	frontEndInput                  string
}

func (suite *CreateOrganizationsSuite) SetupTest() {
	suite.registry = new(mapper.MockRegistryAPI)
	suite.externalUserOrganizationMapper = new(iface.MockExternalUserOrganization)

	suite.input = CreateOrganizationsRequest{
		Registry:         suite.registry,
		Logger:           zaptest.NewLogger(suite.T()),
		OrganizationName: "bui",
	}

	suite.validOrgs = []*model.ExternalUserOrganization{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	suite.validOrgs[0].SetName("kevin")
	suite.validOrgs[1].SetName("charles")
	suite.validOrgs[2].SetName("milu")

	suite.duplicateOrgs = []*model.ExternalUserOrganization{
		{ID: 0},
	}
	suite.duplicateOrgs[0].SetName("bui")
	suite.frontEndInput = "bui"
}

func (suite *CreateOrganizationsSuite) TestCreateOrganization_Success() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Insert", &iface.ExternalUserOrganizationInsertInput{Name: suite.frontEndInput}).Return(&model.ExternalUserOrganization{}, nil)

	_, err := handler(suite.input)
	suite.NoError(err)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Insert", &iface.ExternalUserOrganizationInsertInput{Name: suite.frontEndInput})

}

//this will fail successfully when there is a duplicate
func (suite *CreateOrganizationsSuite) TestCreateOrganization_Fail() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Insert", &iface.ExternalUserOrganizationInsertInput{Name: suite.frontEndInput}).Return(nil, errors.New("ERROR: duplicate key value violates unique constraint \"external_user_organizations_name_key\""))

	_, err := handler(suite.input)
	suite.Error(err)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Insert", &iface.ExternalUserOrganizationInsertInput{Name: suite.frontEndInput})

}

func (suite *CreateOrganizationsSuite) TestErrorResponse() {
	err, _ := errorResponse(errors.New("TEST"), "an error", zaptest.NewLogger(suite.T()))

	suite.NotNil(err)
}

func TestExternalSignUp(t *testing.T) {
	suite.Run(t, new(CreateOrganizationsSuite))
}
