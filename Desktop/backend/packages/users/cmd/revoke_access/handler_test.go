package main

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	cMocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/idp/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ValidateUsernameTestSuite struct {
	suite.Suite
}

func (suite *ValidateUsernameTestSuite) TestValidateUsername_ShouldReturnNilIfUsernameIsValid() {
	input := "Valid Username"
	actual := validateUsername(input)
	suite.Nil(actual)
}

func (suite *ValidateUsernameTestSuite) TestValidUsername_ShouldReturnErrorIfUsernameIsInvalid() {
	input := ""
	actual := validateUsername(input)
	suite.Error(actual)
}

type SendRevokedNotificationTestSuite struct {
	suite.Suite
}

func (suite *SendRevokedNotificationTestSuite) TestSendRevokedNotification_ShouldReturnNilIfCallSendRevokedNotification() {
	const username = "1"
	mockUserDB := new(dynamo.MockDatabase)
	mockExternalWebsocketAPI := new(cMocks.ApiGatewayManagementApiAPI)
	mockLogger := zaptest.NewLogger(suite.T())

	dbExpected := &dynamodb.QueryOutput{
		ConsumedCapacity: &dynamodb.ConsumedCapacity{},
		Count:            new(int64),
		Items:            []map[string]*dynamodb.AttributeValue{},
		LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
		ScannedCount:     new(int64),
	}

	apiExpected := &apigatewaymanagementapi.PostToConnectionOutput{}

	mockUserDB.On("Query", mock.Anything).Return(dbExpected, nil)
	mockExternalWebsocketAPI.On("PostToConnection", mock.Anything).Return(apiExpected, nil)

	actualErr := sendRevokedNotification(&SendRevokedNotificationInput{
		UserID:               username,
		Logger:               mockLogger,
		UsersDB:              mockUserDB,
		ExternalWebsocketAPI: mockExternalWebsocketAPI,
	})

	suite.Nil(actualErr)
}

func (suite *SendRevokedNotificationTestSuite) TestSendRevokedNotification_ShouldReturnErrorIfCallSendRevokedNotificationErrorQuery() {
	const username = "1"
	mockUserDB := new(dynamo.MockDatabase)
	mockExternalWebsocketAPI := new(cMocks.ApiGatewayManagementApiAPI)
	mockLogger := zaptest.NewLogger(suite.T())

	dbExpected := &dynamodb.QueryOutput{
		ConsumedCapacity: &dynamodb.ConsumedCapacity{},
		Count:            new(int64),
		Items:            []map[string]*dynamodb.AttributeValue{},
		LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
		ScannedCount:     new(int64),
	}

	mockUserDB.On("Query", mock.Anything).Return(dbExpected, exception.NewSonarError(http.StatusBadRequest, "test"))

	actualErr := sendRevokedNotification(&SendRevokedNotificationInput{
		UserID:               username,
		Logger:               mockLogger,
		UsersDB:              mockUserDB,
		ExternalWebsocketAPI: mockExternalWebsocketAPI,
	})

	suite.NotNil(actualErr)
	mockUserDB.AssertCalled(suite.T(), "Query", mock.Anything)
	mockExternalWebsocketAPI.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

type HandlerTestSuite struct {
	suite.Suite
}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnNilIfCallHandler() {
	const username = "1"
	const groupname = "test"
	groupnamePointer := new(string)
	*groupnamePointer = groupname
	enabledPointer := new(bool)
	*enabledPointer = true
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}
	userExpected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              enabledPointer,
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}
	group := &cognitoidentityprovider.GroupType{
		CreationDate:     &time.Time{},
		Description:      new(string),
		GroupName:        groupnamePointer,
		LastModifiedDate: &time.Time{},
		Precedence:       new(int64),
		RoleArn:          new(string),
		UserPoolId:       new(string),
	}
	listGroupExpected := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    []*cognitoidentityprovider.GroupType{group},
		NextToken: new(string),
	}

	mockSonarIDP.On("AdminGetUser", username).Return(userExpected, nil)
	mockSonarIDP.On("AdminListGroupsForUser", username).Return(listGroupExpected, nil)
	mockSonarIDP.On("AdminRemoveUserFromGroup", username, groupname).Return(nil, nil)
	mockSonarIDP.On("AdminUserGlobalSignOut", username).Return(nil, nil)
	mockSonarIDP.On("AdminDisableUser", username).Return(nil, nil)

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.Nil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", username, groupname)
	mockSonarIDP.AssertCalled(suite.T(), "AdminUserGlobalSignOut", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminDisableUser", username)
}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnErrorIfCallHandlerAdminGetUserError() {
	const username = "1"
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}

	mockSonarIDP.On("AdminGetUser", mock.Anything).Return(nil, exception.NewSonarError(http.StatusBadRequest, "test"))

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.NotNil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminListGroupsForUser", mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminRemoveUserFromGroup", mock.Anything, mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminUserGlobalSignOut", mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminDisableUser", mock.Anything)

}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnErrorIfCallHandlerAdminListGroupsForUserError() {
	const username = "1"
	enabledPointer := new(bool)
	*enabledPointer = true
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}
	userExpected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              enabledPointer,
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}

	mockSonarIDP.On("AdminGetUser", username).Return(userExpected, nil)
	mockSonarIDP.On("AdminListGroupsForUser", username).Return(nil, exception.NewSonarError(http.StatusBadRequest, "test"))

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.NotNil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", username)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminRemoveUserFromGroup", mock.Anything, mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminUserGlobalSignOut", mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminDisableUser", mock.Anything)

}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnErrorIfCallHandlerAdminRemoveUserFromGroupForUserError() {
	const username = "1"
	const groupname = "test"
	groupnamePointer := new(string)
	*groupnamePointer = groupname
	enabledPointer := new(bool)
	*enabledPointer = true
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}
	userExpected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              enabledPointer,
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}
	group := &cognitoidentityprovider.GroupType{
		CreationDate:     &time.Time{},
		Description:      new(string),
		GroupName:        groupnamePointer,
		LastModifiedDate: &time.Time{},
		Precedence:       new(int64),
		RoleArn:          new(string),
		UserPoolId:       new(string),
	}
	listGroupExpected := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    []*cognitoidentityprovider.GroupType{group},
		NextToken: new(string),
	}

	mockSonarIDP.On("AdminGetUser", username).Return(userExpected, nil)
	mockSonarIDP.On("AdminListGroupsForUser", username).Return(listGroupExpected, nil)
	mockSonarIDP.On("AdminRemoveUserFromGroup", username, groupname).Return(nil, errors.New("test"))

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.NotNil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", username, groupname)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminUserGlobalSignOut", mock.Anything)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminDisableUser", mock.Anything)

}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnErrorIfCallHandlerAdminUserGlobalSignOutError() {
	const username = "1"
	const groupname = "test"
	groupnamePointer := new(string)
	*groupnamePointer = groupname
	enabledPointer := new(bool)
	*enabledPointer = true
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}
	userExpected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              enabledPointer,
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}
	group := &cognitoidentityprovider.GroupType{
		CreationDate:     &time.Time{},
		Description:      new(string),
		GroupName:        groupnamePointer,
		LastModifiedDate: &time.Time{},
		Precedence:       new(int64),
		RoleArn:          new(string),
		UserPoolId:       new(string),
	}
	listGroupExpected := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    []*cognitoidentityprovider.GroupType{group},
		NextToken: new(string),
	}

	mockSonarIDP.On("AdminGetUser", username).Return(userExpected, nil)
	mockSonarIDP.On("AdminListGroupsForUser", username).Return(listGroupExpected, nil)
	mockSonarIDP.On("AdminRemoveUserFromGroup", username, groupname).Return(nil, nil)
	mockSonarIDP.On("AdminUserGlobalSignOut", username).Return(nil, errors.New("test"))

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.NotNil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", username, groupname)
	mockSonarIDP.AssertCalled(suite.T(), "AdminUserGlobalSignOut", username)
	mockSonarIDP.AssertNotCalled(suite.T(), "AdminDisableUser", username)

}

func (suite *HandlerTestSuite) TestHandler_ShouldReturnErrorIfCallHandlerAdminDisableUserError() {
	const username = "1"
	const groupname = "test"
	groupnamePointer := new(string)
	*groupnamePointer = groupname
	enabledPointer := new(bool)
	*enabledPointer = true
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockLogger := zaptest.NewLogger(suite.T())
	mockRequest := &request.RevokeAccessRequest{
		Username: username,
	}
	userExpected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              enabledPointer,
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}
	group := &cognitoidentityprovider.GroupType{
		CreationDate:     &time.Time{},
		Description:      new(string),
		GroupName:        groupnamePointer,
		LastModifiedDate: &time.Time{},
		Precedence:       new(int64),
		RoleArn:          new(string),
		UserPoolId:       new(string),
	}
	listGroupExpected := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    []*cognitoidentityprovider.GroupType{group},
		NextToken: new(string),
	}

	mockSonarIDP.On("AdminGetUser", username).Return(userExpected, nil)
	mockSonarIDP.On("AdminListGroupsForUser", username).Return(listGroupExpected, nil)
	mockSonarIDP.On("AdminRemoveUserFromGroup", username, groupname).Return(nil, nil)
	mockSonarIDP.On("AdminUserGlobalSignOut", username).Return(nil, nil)
	mockSonarIDP.On("AdminDisableUser", username).Return(nil, errors.New("test"))

	actualErr := Handler(&HandlerInput{
		Request:  mockRequest,
		Logger:   mockLogger,
		SonarIDP: mockSonarIDP,
	})

	suite.NotNil(actualErr)
	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", username, groupname)
	mockSonarIDP.AssertCalled(suite.T(), "AdminUserGlobalSignOut", username)
	mockSonarIDP.AssertCalled(suite.T(), "AdminDisableUser", username)

}

/* Execute Suites */
func TestValidateUsernameTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateUsernameTestSuite))
}

func TestSendRevokedNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(SendRevokedNotificationTestSuite))
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
