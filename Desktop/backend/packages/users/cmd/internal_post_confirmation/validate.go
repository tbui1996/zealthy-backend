package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

type ValidateEventInput struct {
	ParsedEvent ParsedEvent
	Idp         cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Environment string
	Username    string
	PoolID      string
}

type ValidatedEvent struct {
	OktaGroup     string
	CognitoGroups []string
}

func validate(input ValidateEventInput) (*ValidatedEvent, *exception.SonarError) {
	// no groups found for environment
	if input.ParsedEvent.OktaSonarGroups == nil || len(input.ParsedEvent.OktaSonarGroups[input.Environment]) == 0 {
		return nil, exception.NewSonarError(http.StatusUnauthorized, "user is trying to log into sonar web, but is not assigned to a group in okta")
	}

	// more than one group found
	if len(input.ParsedEvent.OktaSonarGroups[input.Environment]) > 1 {
		return nil, exception.NewSonarError(http.StatusUnprocessableEntity, "user is trying to log into sonar web, but has more than one group assigned")
	}

	// exactly one sonar group found in okta for an environment, does it exist in cognito??
	oktaGroup := input.ParsedEvent.OktaSonarGroups[input.Environment][0]

	// exists in cognito, so get user's current sonar groups (e.g. prefixed with 'internals_')
	userGroupsOutput, err := input.Idp.AdminListGroupsForUser(&cognitoidentityprovider.AdminListGroupsForUserInput{
		Username:   &input.Username,
		UserPoolId: &input.PoolID,
	})
	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "could not retrieve groups for user")
	}

	// find the congito sonar groups, should only ever be one, but handle case that there is more than one
	var groups []string
	for _, g := range userGroupsOutput.Groups {
		if strings.HasPrefix(*g.GroupName, "internals_") {
			groups = append(groups, *g.GroupName)
		}
	}

	// exactly one group from okta for environment and the mirroring group exists in cognito
	return &ValidatedEvent{
		CognitoGroups: groups,
		OktaGroup:     oktaGroup,
	}, nil
}
