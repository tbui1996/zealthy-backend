package iface

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

type SonarIdentityProvider interface {
	GetListUserInput(params map[string]string) *cognitoidentityprovider.ListUsersInput
	ListUsers(listUserInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error)
	AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	AdminRemoveUserFromGroup(username string, groupName string) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
	AdminUserGlobalSignOut(username string) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error)
	AdminAddUserToGroup(username string, groupName string) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminDisableUser(username string) (*cognitoidentityprovider.AdminDisableUserOutput, error)
	AdminEnableUser(username string) (*cognitoidentityprovider.AdminEnableUserOutput, error)
	AdminListGroupsForUser(username string) (*cognitoidentityprovider.AdminListGroupsForUserOutput, *exception.SonarError)
	GetGroup(group string) (*cognitoidentityprovider.GetGroupOutput, error)
	ListGroups(limit *int64, paginationToken *string) (*cognitoidentityprovider.ListGroupsOutput, error)
}
