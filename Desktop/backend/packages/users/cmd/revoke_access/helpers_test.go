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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type DisableUserTestSuite struct {
	suite.Suite
}

func (suite *DisableUserTestSuite) TestDisableUser_ShouldReturnNilIfCallAdminDisableUserWithUserID() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminDisableUser", mock.Anything).Return(nil, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual := disableUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminDisableUser", userID)
	suite.Nil(actual)
}

func (suite *DisableUserTestSuite) TestDisableUser_ShouldReturnASonarErrorIfAdminDisableUserErrors() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminDisableUser", mock.Anything).Return(nil, errors.New("test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actualErr := disableUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	suite.NotNil(actualErr)
	suite.Error(actualErr)
}

type SignOutUserTestSuite struct {
	suite.Suite
}

func (suite *DisableUserTestSuite) TestSignOutUser_ShouldReturnNilIfCallAdminUserGlobalSignOutWithUserID() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminUserGlobalSignOut", mock.Anything).Return(nil, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual := signOutUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminUserGlobalSignOut", userID)
	suite.Nil(actual)
}

func (suite *DisableUserTestSuite) TestSignOutUser_ShouldReturnASonarErrorIfCallAdminDUserGlobalSignOutrErrors() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminUserGlobalSignOut", mock.Anything).Return(nil, errors.New("test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actualErr := signOutUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	suite.NotNil(actualErr)
	suite.Error(actualErr)
}

type ListGroupsTestSuite struct {
	suite.Suite
}

func (suite *ListGroupsTestSuite) TestListGroup_ShouldReturnListIfCallAdminListGroupsForUserWithUserID() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	expected := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups:    []*cognitoidentityprovider.GroupType{},
		NextToken: new(string),
	}
	mockSonarIDP.On("AdminListGroupsForUser", mock.Anything).Return(expected, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual, actualErr := listGroups(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", userID)
	suite.NotNil(actual)
	suite.Nil(actualErr)
	suite.Equal(expected, actual)
}

func (suite *ListGroupsTestSuite) TestListGroup_ShouldReturnASonarErrorIfCallAdminListGroupsForUserErrors() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminListGroupsForUser", mock.Anything).Return(nil, exception.NewSonarError(http.StatusBadRequest, "test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual, actualErr := listGroups(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminListGroupsForUser", userID)
	suite.Nil(actual)
	suite.Error(actualErr)
}

type RemoveUsersFromGroupTestSuite struct {
	suite.Suite
}

func (suite *RemoveUsersFromGroupTestSuite) TestRemoveUserFromGroup_ShouldReturnNilIfCallAdminRemoveUserFromGroup() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminRemoveUserFromGroup", mock.Anything, mock.Anything).Return(nil, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	groupName := "GroupTest"

	actual := removeUserFromGroup(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	}, groupName)

	mockSonarIDP.AssertCalled(suite.T(), "AdminRemoveUserFromGroup", userID, groupName)
	suite.Nil(actual)
}

func (suite *DisableUserTestSuite) TestRemoveUserFromGroup_ShouldReturnASonarErrorIfCallRemoveUserFromGroupErrors() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminRemoveUserFromGroup", mock.Anything, mock.Anything).Return(nil, errors.New("test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	groupName := "GroupTest"

	actualErr := removeUserFromGroup(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	}, groupName)

	suite.NotNil(actualErr)
	suite.Error(actualErr)
}

type GetUserTestSuite struct {
	suite.Suite
}

func (suite *GetUserTestSuite) TestGetUser_ShouldReturnUserIfCallAdminGetUser() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	expected := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled:              new(bool),
		MFAOptions:           []*cognitoidentityprovider.MFAOptionType{},
		PreferredMfaSetting:  new(string),
		UserAttributes:       []*cognitoidentityprovider.AttributeType{},
		UserCreateDate:       &time.Time{},
		UserLastModifiedDate: &time.Time{},
		UserMFASettingList:   []*string{},
		UserStatus:           new(string),
		Username:             new(string),
	}

	mockSonarIDP.On("AdminGetUser", mock.Anything).Return(expected, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual, actualErr := getUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", userID)
	suite.NotNil(actual)
	suite.Nil(actualErr)
	suite.Equal(expected, actual)
}

func (suite *GetUserTestSuite) TestGetUser_ShouldReturnASonarErrorIfCallAdminGetUserErrors() {
	mockSonarIDP := new(mocks.SonarIdentityProvider)
	mockSonarIDP.On("AdminGetUser", mock.Anything).Return(nil, exception.NewSonarError(http.StatusBadRequest, "test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"
	actual, actualErr := getUser(&HandleRevokeUserInput{
		SonarIDP: mockSonarIDP,
		Logger:   mockLogger,
		UserID:   userID,
	})

	mockSonarIDP.AssertCalled(suite.T(), "AdminGetUser", userID)
	suite.Nil(actual)
	suite.Error(actualErr)
}

type GetSonarWebsocketConnectionTestSuite struct {
	suite.Suite
}

func (suite *GetSonarWebsocketConnectionTestSuite) TestGetSonarWebsocketConnection_ShouldReturnQueryIfCallGetSonarWebsocketConnection() {
	mockUserDB := new(dynamo.MockDatabase)
	expected := &dynamodb.QueryOutput{
		ConsumedCapacity: &dynamodb.ConsumedCapacity{},
		Count:            new(int64),
		Items:            []map[string]*dynamodb.AttributeValue{},
		LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
		ScannedCount:     new(int64),
	}

	mockUserDB.On("Query", mock.Anything).Return(expected, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"

	actual, actualErr := getSonarWebsocketConnection(&GetSonarWebsocketConnectionInput{
		UserID:  userID,
		Logger:  mockLogger,
		UsersDB: mockUserDB,
	})

	mockUserDB.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.NotNil(actual)
	suite.Nil(actualErr)
	suite.Equal(expected, actual)
}

func (suite *GetSonarWebsocketConnectionTestSuite) TestGetSonarWebsocketConnection_ShouldReturnNilIfCallGetSonarWebsocketConnectionError() {
	mockUserDB := new(dynamo.MockDatabase)
	output := &dynamodb.QueryOutput{
		ConsumedCapacity: &dynamodb.ConsumedCapacity{},
		Count:            new(int64),
		Items:            []map[string]*dynamodb.AttributeValue{},
		LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
		ScannedCount:     new(int64),
	}

	mockUserDB.On("Query", mock.Anything).Return(output, exception.NewSonarError(http.StatusBadRequest, "test"))

	mockLogger := zaptest.NewLogger(suite.T())

	userID := "1"

	actual, actualErr := getSonarWebsocketConnection(&GetSonarWebsocketConnectionInput{
		UserID:  userID,
		Logger:  mockLogger,
		UsersDB: mockUserDB,
	})

	mockUserDB.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.Nil(actual)
	suite.NotNil(actualErr)
	suite.Error(actualErr)
}

type PostConnectionTestSuite struct {
	suite.Suite
}

func (suite *PostConnectionTestSuite) TestPostConnection_ShouldReturnNillIfCallPostConnection() {
	mockExternalWebsocketAPI := new(cMocks.ApiGatewayManagementApiAPI)

	expected := &apigatewaymanagementapi.PostToConnectionOutput{}

	mockExternalWebsocketAPI.On("PostToConnection", mock.Anything).Return(expected, nil)

	mockLogger := zaptest.NewLogger(suite.T())

	jsonData := []byte("Here is a string....")

	actualErr := postConnection(&PostConnectionInput{
		ConnectionID:         nil,
		Data:                 jsonData,
		Logger:               mockLogger,
		ExternalWebsocketAPI: mockExternalWebsocketAPI,
	})

	mockExternalWebsocketAPI.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	suite.Nil(actualErr)
}

func (suite *PostConnectionTestSuite) TestPostConnection_ShouldReturnErrorIfCallPostConnectionError() {
	mockExternalWebsocketAPI := new(cMocks.ApiGatewayManagementApiAPI)

	mockExternalWebsocketAPI.On("PostToConnection", mock.Anything).Return(nil, errors.New("test"))

	mockLogger := zaptest.NewLogger(suite.T())

	jsonData := []byte("Here is a string....")

	actualErr := postConnection(&PostConnectionInput{
		ConnectionID:         nil,
		Data:                 jsonData,
		Logger:               mockLogger,
		ExternalWebsocketAPI: mockExternalWebsocketAPI,
	})

	mockExternalWebsocketAPI.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	suite.NotNil(actualErr)
}

/* Execute Suites */

func TestDisableUserTestSuite(t *testing.T) {
	suite.Run(t, new(DisableUserTestSuite))
}

func TestHandleDisableAndSignOutTestSuite(t *testing.T) {
	suite.Run(t, new(SignOutUserTestSuite))
}

func TestListGroupsTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(ListGroupsTestSuite))
}

func TestRemoveUsersFromGroupTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveUsersFromGroupTestSuite))
}

func TestGetUserTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserTestSuite))
}

func TestGetSonarWebsocketConnectionTestSuite(t *testing.T) {
	suite.Run(t, new(GetSonarWebsocketConnectionTestSuite))
}

func TestPostConnectionTestSuite(t *testing.T) {
	suite.Run(t, new(PostConnectionTestSuite))
}
