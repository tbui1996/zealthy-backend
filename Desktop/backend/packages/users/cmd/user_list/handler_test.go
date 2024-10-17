package main

import (
	"errors"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/users/pkg/dto"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
	input              HandlerInput
	externalUserMapper *iface.MockExternalUser
	registry           *mapper.MockRegistryAPI
}

func (suite *HandlerSuite) SetupTest() {
	suite.registry = new(mapper.MockRegistryAPI)
	suite.externalUserMapper = new(iface.MockExternalUser)

	suite.input = HandlerInput{
		QueryStringParameters: map[string]string{},
		Registry:              suite.registry,
		Logger:                zaptest.NewLogger(suite.T()),
	}
}

func (suite *HandlerSuite) TestHandlerErrGetListUserInput() {
	suite.registry.On("ExternalUser").Return(suite.externalUserMapper)
	suite.externalUserMapper.On("FindAll").Return(nil, errors.New("test"))

	actualOut, actualErr := handler(suite.input)
	suite.Nil(actualOut)
	suite.Error(actualErr)
	suite.externalUserMapper.AssertCalled(suite.T(), "FindAll")
}

func (suite *HandlerSuite) TestHandlerSuccess() {
	organization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).WithName("testOrganization").Value()

	user1 := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test@gmail.com",
		Username: "test@gmail.com",
		Email:    "test@gmail.com",
		Status:   "enabled",
	}).WithEnabled(true).WithFirstName("firstName").WithLastName("lastName").WithGroup("group").WithOrganization(organization).Value()

	user2 := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test2@gmail.com",
		Username: "test2@gmail.com",
		Email:    "test2@gmail.com",
		Status:   "enabled",
	}).Value()

	output := []*model.ExternalUser{
		user1,
		user2,
	}

	suite.registry.On("ExternalUser").Return(suite.externalUserMapper)
	suite.externalUserMapper.On("FindAll").Return(output, nil)

	actualOut, actualErr := handler(suite.input)
	suite.Nil(actualErr)
	suite.externalUserMapper.AssertCalled(suite.T(), "FindAll")

	expected := []*dto.ExternalUser{
		dto.ExternalUserFromModel(user1),
		dto.ExternalUserFromModel(user2),
	}

	suite.Len(actualOut, 2)
	suite.Equal(expected, actualOut)
}

func TestExternalSignUp(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
