package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
	groupName   string
	environment string
	poolID      string
	username    string
}

func (suite *HandlerSuite) SetupSuite() {
	suite.groupName = "internals_program_manager"
	suite.environment = "dev"
	suite.poolID = "pool_id"
	suite.username = "Okta_1234567890"
}

func (suite *HandlerSuite) TestHandler_OneCognitoGroupAndItMatchesOkta() {
	input := HandlerInput{
		Idp:      new(mocks.CognitoIdentityProviderAPI),
		PoolID:   suite.poolID,
		Username: suite.username,
		ValidatedEvent: &ValidatedEvent{
			OktaGroup:     suite.groupName,
			CognitoGroups: []string{suite.groupName},
		},
		Logger: zaptest.NewLogger(suite.T()),
	}

	actualErr := handler(input)
	suite.Nil(actualErr)
}

func (suite *HandlerSuite) TestHandler_ErrRemovingCurrentGroups() {
	idpMock := new(mocks.CognitoIdentityProviderAPI)
	adminRemoveUserFromGroup := &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		GroupName:  &suite.groupName,
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}
	idpMock.On("AdminRemoveUserFromGroup", adminRemoveUserFromGroup).Return(nil, errors.New("idp err"))

	input := HandlerInput{
		Idp:      idpMock,
		PoolID:   suite.poolID,
		Username: suite.username,
		ValidatedEvent: &ValidatedEvent{
			OktaGroup:     "different_group",
			CognitoGroups: []string{suite.groupName},
		},
		Logger: zaptest.NewLogger(suite.T()),
	}

	actualErr := handler(input)
	suite.Equal(http.StatusBadRequest, actualErr.StatusCode)
	idpMock.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", adminRemoveUserFromGroup)
}

func (suite *HandlerSuite) TestHandler_NoCurrentGroups() {
	idpMock := new(mocks.CognitoIdentityProviderAPI)
	addUserToGroupInput := &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &suite.groupName,
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}
	idpMock.On("AdminAddUserToGroup", addUserToGroupInput).Return(nil, nil)

	input := HandlerInput{
		Idp:      idpMock,
		PoolID:   suite.poolID,
		Username: suite.username,
		ValidatedEvent: &ValidatedEvent{
			OktaGroup:     suite.groupName,
			CognitoGroups: []string{},
		},
		Logger: zaptest.NewLogger(suite.T()),
	}

	actualErr := handler(input)
	suite.Nil(actualErr)
	idpMock.AssertCalled(suite.T(), "AdminAddUserToGroup", addUserToGroupInput)
}

func (suite *HandlerSuite) TestHandler_ErrAssignGroups() {
	idpMock := new(mocks.CognitoIdentityProviderAPI)
	addUserToGroupInput := &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &suite.groupName,
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}
	idpMock.On("AdminAddUserToGroup", addUserToGroupInput).Return(nil, errors.New("idp err"))

	input := HandlerInput{
		Idp:      idpMock,
		PoolID:   suite.poolID,
		Username: suite.username,
		ValidatedEvent: &ValidatedEvent{
			OktaGroup:     suite.groupName,
			CognitoGroups: []string{},
		},
		Logger: zaptest.NewLogger(suite.T()),
	}

	actualErr := handler(input)
	suite.Equal(http.StatusBadRequest, actualErr.StatusCode)
	idpMock.AssertCalled(suite.T(), "AdminAddUserToGroup", addUserToGroupInput)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
