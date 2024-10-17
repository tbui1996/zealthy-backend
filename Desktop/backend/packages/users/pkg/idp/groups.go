package idp

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func (idp *CognitoSonarIdentityProvider) AdminListGroupsForUser(username string) (*cognitoidentityprovider.AdminListGroupsForUserOutput, *exception.SonarError) {

	groupOutput, err := idp.svc.AdminListGroupsForUser(&cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(idp.userPoolID),
	})
	if err != nil {
		return nil, HandleAwsError(err)
	}

	return groupOutput, nil
}

func (idp *CognitoSonarIdentityProvider) GetGroup(group string) (*cognitoidentityprovider.GetGroupOutput, error) {
	return idp.svc.GetGroup(&cognitoidentityprovider.GetGroupInput{
		GroupName:  aws.String(group),
		UserPoolId: aws.String(idp.userPoolID),
	})
}

func (idp *CognitoSonarIdentityProvider) ListGroups(limit *int64, paginationToken *string) (*cognitoidentityprovider.ListGroupsOutput, error) {
	var input cognitoidentityprovider.ListGroupsInput
	input.UserPoolId = aws.String(idp.userPoolID)

	if limit != nil {
		input.Limit = limit
	} else {
		var maxLimit int64 = 60
		input.Limit = &maxLimit
	}

	if paginationToken != nil {
		input.NextToken = paginationToken
	}

	return idp.svc.ListGroups(&input)
}
