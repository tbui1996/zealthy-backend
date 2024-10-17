package idp

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (idp *CognitoSonarIdentityProvider) GetListUserInput(params map[string]string) *cognitoidentityprovider.ListUsersInput {
	listUserInput := &cognitoidentityprovider.ListUsersInput{}

	paginationToken := params["paginationToken"]
	if len(paginationToken) > 0 && paginationToken != "null" {
		listUserInput.PaginationToken = aws.String(paginationToken)
	}

	// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUsers.html#API_ListUsers_RequestSyntax
	// if limit is not passed, or cannot be parsed, use max limit of 60
	var maxLimit int64 = 60

	limitStr := params["limit"]
	if len(limitStr) > 0 {
		limit, err := strconv.ParseInt(limitStr, 10, 64) //nolint
		if err == nil && limit < maxLimit {
			listUserInput.Limit = aws.Int64(limit)
		} else {
			listUserInput.Limit = aws.Int64(maxLimit)
		}
	}

	listUserInput.UserPoolId = aws.String(idp.userPoolID)

	return listUserInput
}

func (idp *CognitoSonarIdentityProvider) ListUsers(listUserInput *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	return idp.svc.ListUsers(listUserInput)
}

func (idp *CognitoSonarIdentityProvider) AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	return idp.svc.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminRemoveUserFromGroup(username string, groupName string) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
	return idp.svc.AdminRemoveUserFromGroup(&cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		Username:   aws.String(username),
		GroupName:  aws.String(groupName),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminUserGlobalSignOut(username string) (*cognitoidentityprovider.AdminUserGlobalSignOutOutput, error) {
	return idp.svc.AdminUserGlobalSignOut(&cognitoidentityprovider.AdminUserGlobalSignOutInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminAddUserToGroup(username string, groupName string) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	return idp.svc.AdminAddUserToGroup(&cognitoidentityprovider.AdminAddUserToGroupInput{
		Username:   aws.String(username),
		GroupName:  aws.String(groupName),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminDisableUser(username string) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	return idp.svc.AdminDisableUser(&cognitoidentityprovider.AdminDisableUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminEnableUser(username string) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	return idp.svc.AdminEnableUser(&cognitoidentityprovider.AdminEnableUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) AdminInitiateAuth(username string, clientId string) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
	return idp.svc.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
		},
		ClientId:   aws.String(clientId),
		UserPoolId: aws.String(idp.userPoolID),
	})
}
