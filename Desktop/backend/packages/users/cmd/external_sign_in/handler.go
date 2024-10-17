package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/response"
	"go.uber.org/zap"
)

type HandlerInput struct {
	Logger   *zap.Logger
	Idp      cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Email    string
	ClientID string
	PoolID   string
}

func handler(input HandlerInput) (string, *exception.SonarError) {
	input.Logger.Debug("getting user")

	user, err := input.Idp.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(input.PoolID),
		Username:   aws.String(input.Email),
	})
	if err != nil {
		return "", handleAdminGetUserErr(err)
	}
	if user == nil {
		return "", exception.NewSonarError(http.StatusInternalServerError, "expected user to not be null")
	}
	if !*user.Enabled {
		return "", exception.NewSonarError(http.StatusUnauthorized, "user is registered, awaiting confirmation and role assignment")
	}
	if *user.UserStatus != cognitoidentityprovider.UserStatusTypeConfirmed {
		input.Logger.Debug("unknown user status: " + *user.UserStatus)

		return "", exception.NewSonarError(http.StatusInternalServerError, "unknown user status: "+*user.UserStatus)
	}

	input.Logger.Debug("handling signed up and confirmed user")

	initiateAuth, err := input.Idp.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(input.Email),
		},
		ClientId:   aws.String(input.ClientID),
		UserPoolId: aws.String(input.PoolID),
	})
	if err != nil {
		errMessage := fmt.Sprintf("error initiating auth (%s)", err)
		return "", exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	authPayload := response.AuthResponsePayload{
		ID:           *initiateAuth.AuthenticationResult.IdToken,
		Token:        *initiateAuth.AuthenticationResult.AccessToken,
		RefreshToken: *initiateAuth.AuthenticationResult.RefreshToken,
	}

	jsonPayload, err := json.Marshal(authPayload)
	if err != nil {
		return "", exception.NewSonarError(http.StatusInternalServerError, err.Error())
	}

	input.Logger.Debug("Signed In")

	return string(jsonPayload), nil
}
