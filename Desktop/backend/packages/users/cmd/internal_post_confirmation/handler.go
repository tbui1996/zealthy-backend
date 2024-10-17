package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"go.uber.org/zap"
)

type HandlerInput struct {
	Idp            cognitoidentityprovideriface.CognitoIdentityProviderAPI
	PoolID         string
	Username       string
	ValidatedEvent *ValidatedEvent
	Logger         *zap.Logger
}

func handler(input HandlerInput) *exception.SonarError {
	// validated event contains one and only one okta group that the user should be assigned, and all current cognito 'sonar' groups assigned

	if len(input.ValidatedEvent.CognitoGroups) == 1 && input.ValidatedEvent.CognitoGroups[0] == input.ValidatedEvent.OktaGroup {
		input.Logger.Debug("user has one group assigned in cognito, and that group matches okta")
		return nil
	}

	// remove current 'sonar' groups if there are any
	var errors []string
	for _, group := range input.ValidatedEvent.CognitoGroups {
		// pin group so tests don't lose reference to group that was called
		group := group

		_, err := input.Idp.AdminRemoveUserFromGroup(&cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			GroupName:  &group,
			Username:   &input.Username,
			UserPoolId: &input.PoolID,
		})
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		input.Logger.Error("could not remove groups from user" + strings.Join(errors, ", "))
		return exception.NewSonarError(http.StatusBadRequest, "could not remove groups from user"+strings.Join(errors, ", "))
	}

	// assign cognito group to user
	_, err := input.Idp.AdminAddUserToGroup(&cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  &input.ValidatedEvent.OktaGroup,
		Username:   &input.Username,
		UserPoolId: &input.PoolID,
	})

	if err != nil {
		input.Logger.Error("could not assign user to group: " + err.Error())
		return exception.NewSonarError(http.StatusBadRequest, "could not assign user to group: "+err.Error())
	}

	return nil
}
