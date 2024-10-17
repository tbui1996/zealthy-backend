package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/suite"
)

type ValidationSuite struct {
	suite.Suite
	groupName   string
	environment string
	poolID      string
	username    string
}

func (suite *ValidationSuite) SetupSuite() {
	suite.groupName = "internals_program_manager"
	suite.environment = "dev"
	suite.poolID = "pool_id"
	suite.username = "Okta_1234567890"
}

func (suite *ValidationSuite) TestValidate_OneGroup() {
	groupsMap := map[string][]string{"dev": {suite.groupName}}

	idpMock := new(mocks.CognitoIdentityProviderAPI)
	idpMock.On("AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}).Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{{GroupName: &suite.groupName}},
	}, nil)

	expectedValidatedEvent := &ValidatedEvent{OktaGroup: suite.groupName, CognitoGroups: []string{suite.groupName}}

	actualValidatedEvent, actualErr := validate(ValidateEventInput{
		ParsedEvent: ParsedEvent{
			OktaSonarGroups: groupsMap,
		},
		Idp:         idpMock,
		Environment: suite.environment,
		PoolID:      suite.poolID,
		Username:    suite.username,
	})

	suite.Equal(expectedValidatedEvent, actualValidatedEvent)
	suite.Nil(actualErr)
	idpMock.AssertCalled(suite.T(), "AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	})
}

func (suite *ValidationSuite) TestValidate_MoreThanOneGroupButForDifferentEnvironments() {
	groupsMap := map[string][]string{"dev": {suite.groupName}, "prod": {suite.groupName}}

	idpMock := new(mocks.CognitoIdentityProviderAPI)
	idpMock.On("AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}).Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{{GroupName: &suite.groupName}},
	}, nil)

	expectedValidatedEvent := &ValidatedEvent{OktaGroup: suite.groupName, CognitoGroups: []string{suite.groupName}}

	actualValidatedEvent, actualErr := validate(ValidateEventInput{
		ParsedEvent: ParsedEvent{
			OktaSonarGroups: groupsMap,
		},
		Idp:         idpMock,
		Environment: suite.environment,
		PoolID:      suite.poolID,
		Username:    suite.username,
	})

	suite.Equal(expectedValidatedEvent, actualValidatedEvent)
	suite.Nil(actualErr)
	idpMock.AssertCalled(suite.T(), "AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	})
}

func (suite *ValidationSuite) TestValidate_NoSonarGroup() {
	groupsMap := map[string][]string{}

	idpMock := new(mocks.CognitoIdentityProviderAPI)

	actualValidatedEvent, actualErr := validate(ValidateEventInput{
		ParsedEvent: ParsedEvent{
			OktaSonarGroups: groupsMap,
		},
		Idp:         idpMock,
		Environment: suite.environment,
		PoolID:      suite.poolID,
		Username:    suite.username,
	})

	suite.Nil(actualValidatedEvent)
	suite.Equal(http.StatusUnauthorized, actualErr.StatusCode)
	idpMock.AssertNotCalled(suite.T(), "AdminListGroupsForUser")
}

func (suite *ValidationSuite) TestValidate_MoreThanOneGroupForAnEnvironment() {
	groupsMap := map[string][]string{"dev": {"internals_program_manager", "internals_NOT_program_manager"}}

	idpMock := new(mocks.CognitoIdentityProviderAPI)

	actualValidatedEvent, actualErr := validate(ValidateEventInput{
		ParsedEvent: ParsedEvent{
			OktaSonarGroups: groupsMap,
		},
		Idp:         idpMock,
		Environment: suite.environment,
		PoolID:      suite.poolID,
		Username:    suite.username,
	})

	suite.Nil(actualValidatedEvent)
	suite.Equal(http.StatusUnprocessableEntity, actualErr.StatusCode)
	idpMock.AssertNotCalled(suite.T(), "AdminListGroupsForUser")
}

func (suite *ValidationSuite) TestValidate_IdpErr() {
	groupsMap := map[string][]string{"dev": {suite.groupName}}

	idpMock := new(mocks.CognitoIdentityProviderAPI)
	idpMock.On("AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	}).Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{{GroupName: &suite.groupName}},
	}, errors.New("idp err"))

	actualValidatedEvent, actualErr := validate(ValidateEventInput{
		ParsedEvent: ParsedEvent{OktaSonarGroups: groupsMap},
		Idp:         idpMock,
		Environment: suite.environment,
		PoolID:      suite.poolID,
		Username:    suite.username,
	})

	suite.Nil(actualValidatedEvent)
	suite.Equal(http.StatusBadRequest, actualErr.StatusCode)
	idpMock.AssertCalled(suite.T(), "AdminListGroupsForUser", &cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &suite.username,
		UserPoolId: &suite.poolID,
	})
}

func TestValidationSuite(t *testing.T) {
	suite.Run(t, new(ValidationSuite))
}
