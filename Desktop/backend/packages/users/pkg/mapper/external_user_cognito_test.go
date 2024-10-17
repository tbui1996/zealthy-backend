package mapper

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"testing"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ExternalUserCognitoSuite struct {
	suite.Suite
}

func (s *ExternalUserCognitoSuite) TestValidateUserType_ValidUser() {
	enabled := true
	userStatus := "CONFIRMED"
	username := "test"
	email := "test@gmail.com"
	emailName := "email"

	valid := &cognitoidentityprovider.UserType{
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email,
			},
		},
		Enabled:    &enabled,
		UserStatus: &userStatus,
		Username:   &username,
	}

	err := validateUserType(valid)
	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestValidateUserType_ErrorsWithoutValidAttribute() {
	enabled := true
	userStatus := "CONFIRMED"
	username := "test"

	valid := &cognitoidentityprovider.UserType{
		Attributes: []*cognitoidentityprovider.AttributeType{},
		Enabled:    &enabled,
		UserStatus: &userStatus,
		Username:   &username,
	}

	err := validateUserType(valid)
	s.Error(err, "Expected to find email in UserAttributes")
}

func (s *ExternalUserCognitoSuite) TestParseUserTypeToRecord_MapsUserTypeTypeToRecord() {
	enabled := true
	userStatus := "CONFIRMED"
	username := "test"
	email := "test@gmail.com"
	emailName := "email"
	user := &cognitoidentityprovider.UserType{
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email,
			},
		},
		Enabled:    &enabled,
		UserStatus: &userStatus,
		Username:   &username,
	}

	record := &externalUserCognitoRecord{}
	parseUserTypeToRecord(record, user)

	s.Equal(record.email, email)
	s.Equal(record.enabled, enabled)
	s.Equal(record.status, userStatus)
	s.Equal(record.username, username)
}

func (s *ExternalUserCognitoSuite) TestFindCognitoGroupRecord_ReturnsGroupName() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	group := "testGroup"
	output := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group,
			},
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(output, nil)

	groupName, hasGroup, err := m.findCognitoGroupRecord("test")

	s.NoError(err)
	s.True(hasGroup)
	s.Equal(group, groupName)
}

func (s *ExternalUserCognitoSuite) TestFindCognitoGroupRecord_ReturnsNoGroupIfMultipleGroups() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	group := "testGroup"
	output := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group,
			},
			{
				GroupName: &group,
			},
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(output, nil)

	groupName, hasGroup, err := m.findCognitoGroupRecord("test")

	s.NoError(err)
	s.False(hasGroup)
	s.Equal("", groupName)
}

func (s *ExternalUserCognitoSuite) TestFindCognitoGroupRecord_ErrorsIfGroupNameIsNil() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	output := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: nil,
			},
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(output, nil)

	groupName, hasGroup, err := m.findCognitoGroupRecord("test")

	s.Error(err)
	s.False(hasGroup)
	s.Equal("", groupName)
}

func (s *ExternalUserCognitoSuite) TestBuildCompleteRecordFromUserType_ReturnsRecord() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	enabled := true
	userStatus := "CONFIRMED"
	username := "test"
	email := "test@gmail.com"
	emailName := "email"
	user := &cognitoidentityprovider.UserType{
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email,
			},
		},
		Enabled:    &enabled,
		UserStatus: &userStatus,
		Username:   &username,
	}

	group := "testGroup"
	output := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group,
			},
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(output, nil)

	actual, err := m.buildCompleteRecordFromUserType(user)

	expected := &externalUserCognitoRecord{
		username:  username,
		status:    userStatus,
		enabled:   enabled,
		email:     email,
		firstName: nil,
		lastName:  nil,
		group:     group,
		hasGroup:  true,
	}

	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *ExternalUserCognitoSuite) TestFind_ReturnsRecord() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	enabled := true
	userStatus := "CONFIRMED"
	username := "test"
	email := "test@gmail.com"
	emailName := "email"
	group := "testGroup"

	adminListGroupsForUserOutput := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group,
			},
		},
	}

	adminGetUserOutput := &cognitoidentityprovider.AdminGetUserOutput{
		Enabled: &enabled,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email,
			},
		},
		UserStatus: &userStatus,
		Username:   &username,
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(adminListGroupsForUserOutput, nil)
	idp.On("AdminGetUser", mock.Anything).Return(adminGetUserOutput, nil)

	actual, err := m.find("test")

	expected := &externalUserCognitoRecord{
		username:  username,
		status:    userStatus,
		enabled:   enabled,
		email:     email,
		firstName: nil,
		lastName:  nil,
		group:     group,
		hasGroup:  true,
	}

	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *ExternalUserCognitoSuite) TestUpdateAttributes_Success() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	req := externalUserCognitoRecordUpdater{
		firstName: "test",
		lastName:  "person",
		username:  "test_person",
	}

	idp.On("AdminUpdateUserAttributes", mock.Anything).Return(nil, nil)

	err := m.updateAttributes(req)

	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestClearGroups_Success() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	resp := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{GroupName: aws.String("group1")},
			{GroupName: aws.String("group2")},
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(resp, nil)
	idp.On("AdminRemoveUserFromGroup", mock.Anything).Return(nil, nil)

	err := m.clearGroups("1")
	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestUpdateGroup_Success() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	idp.On("AdminAddUserToGroup", mock.Anything).Return(nil, nil)

	err := m.updateGroup("1", "group")
	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestUpdateEnabled_SuccessEnable() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	idp.On("AdminEnableUser", mock.Anything).Return(nil, nil)

	err := m.updateEnabled("1", true)
	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestUpdateEnabled_SuccessDisable() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	idp.On("AdminDisableUser", mock.Anything).Return(nil, nil)

	err := m.updateEnabled("1", false)
	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestUpdate_Success() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	req := externalUserCognitoRecordUpdater{
		firstName:       "test",
		lastName:        "person",
		username:        "test_person",
		enabledChanged:  true,
		enabled:         true,
		hasGroupChanged: true,
		hasGroup:        false,
	}

	resp := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{GroupName: aws.String("group1")},
			{GroupName: aws.String("group2")},
		},
	}

	idp.On("AdminEnableUser", mock.Anything).Return(nil, nil)
	idp.On("AdminListGroupsForUser", mock.Anything).Return(resp, nil)
	idp.On("AdminRemoveUserFromGroup", mock.Anything).Return(nil, nil)

	err := m.updateAttributes(req)

	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestUpdate_SuccessGroupChanged() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	req := externalUserCognitoRecordUpdater{
		firstName:       "test",
		lastName:        "person",
		username:        "test_person",
		enabledChanged:  true,
		enabled:         true,
		hasGroupChanged: false,
		hasGroup:        false,
		groupChanged:    true,
	}

	resp := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{GroupName: aws.String("group1")},
			{GroupName: aws.String("group2")},
		},
	}

	idp.On("AdminEnableUser", mock.Anything).Return(nil, nil)
	idp.On("AdminListGroupsForUser", mock.Anything).Return(resp, nil)
	idp.On("AdminAddUserToGroup", mock.Anything).Return(nil, nil)

	err := m.updateAttributes(req)

	s.NoError(err)
}

func (s *ExternalUserCognitoSuite) TestFindAll_ReturnsRecords() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	emailName := "email"
	group := "testGroup"

	enabled1 := true
	userStatus1 := "CONFIRMED"
	username1 := "test"
	email1 := "test@gmail.com"

	enabled2 := false
	userStatus2 := "UNCONFIRMED"
	username2 := "test2"
	email2 := "test2@gmail.com"

	user1 := &cognitoidentityprovider.UserType{
		Enabled: &enabled1,
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email1,
			},
		},
		UserStatus: &userStatus1,
		Username:   &username1,
	}

	user2 := &cognitoidentityprovider.UserType{
		Enabled: &enabled2,
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email2,
			},
		},
		UserStatus: &userStatus2,
		Username:   &username2,
	}

	adminListGroupsForUserOutput := &cognitoidentityprovider.AdminListGroupsForUserOutput{
		Groups: []*cognitoidentityprovider.GroupType{
			{
				GroupName: &group,
			},
		},
	}

	listUsersOutput := &cognitoidentityprovider.ListUsersOutput{
		Users: []*cognitoidentityprovider.UserType{
			user1,
			user2,
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(adminListGroupsForUserOutput, nil)
	idp.On("ListUsers", mock.Anything).Return(listUsersOutput, nil)

	actual, err := m.findAll()

	expected := []*externalUserCognitoRecord{
		{
			username:  username1,
			status:    userStatus1,
			enabled:   enabled1,
			email:     email1,
			firstName: nil,
			lastName:  nil,
			group:     group,
			hasGroup:  true,
		},
		{
			username:  username2,
			status:    userStatus2,
			enabled:   enabled2,
			email:     email2,
			firstName: nil,
			lastName:  nil,
			group:     group,
			hasGroup:  true,
		},
	}

	s.NoError(err)
	s.Len(actual, 2)
	for _, actualRecord := range actual {
		exists := false

		for _, expectedRecord := range expected {
			if !exists {
				exists = *actualRecord == *expectedRecord
			}
		}

		s.True(exists)
	}
}

func (s *ExternalUserCognitoSuite) TestFindAll_ReturnsConcatenatedErrors() {
	idp := new(mocks.CognitoIdentityProviderAPI)

	m := &externalUserCognito{
		idp:        idp,
		userPoolId: "test",
		logger:     zaptest.NewLogger(s.T()),
	}

	emailName := "email"

	enabled1 := true
	userStatus1 := "CONFIRMED"
	username1 := "test"
	email1 := "test@gmail.com"

	enabled2 := false
	userStatus2 := "UNCONFIRMED"
	username2 := "test2"
	email2 := "test2@gmail.com"

	user1 := &cognitoidentityprovider.UserType{
		Enabled: &enabled1,
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email1,
			},
		},
		UserStatus: &userStatus1,
		Username:   &username1,
	}

	user2 := &cognitoidentityprovider.UserType{
		Enabled: &enabled2,
		Attributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  &emailName,
				Value: &email2,
			},
		},
		UserStatus: &userStatus2,
		Username:   &username2,
	}

	listUsersOutput := &cognitoidentityprovider.ListUsersOutput{
		Users: []*cognitoidentityprovider.UserType{
			user1,
			user2,
		},
	}

	idp.On("AdminListGroupsForUser", mock.Anything).Return(nil, errors.New("error"))
	idp.On("ListUsers", mock.Anything).Return(listUsersOutput, nil)

	actual, err := m.findAll()

	s.Error(err, "error; error")
	s.Nil(actual)
}

func TestExternalUserCognitoSuite(t *testing.T) {
	suite.Run(t, new(ExternalUserCognitoSuite))
}
