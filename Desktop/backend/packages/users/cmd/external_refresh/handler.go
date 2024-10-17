package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/externalAuth"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/idp"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/response"
)

type HandlerInput struct {
	ClientID     string
	Idp          cognitoidentityprovideriface.CognitoIdentityProviderAPI
	RefreshToken string
}

func handler(input HandlerInput) ([]byte, *exception.SonarError) {
	initiateAuth, err := input.Idp.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeRefreshToken),
		ClientId: aws.String(input.ClientID),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(input.RefreshToken),
		},
	})
	if err != nil {
		return nil, idp.HandleAwsError(err)
	}

	err = externalAuth.ValidateInitiateAuthOutput(initiateAuth)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, err.Error())
	}

	authPayload := response.AuthResponsePayload{
		ID:    *initiateAuth.AuthenticationResult.IdToken,
		Token: *initiateAuth.AuthenticationResult.AccessToken,
	}

	jsonPayload, err := json.Marshal(authPayload)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, fmt.Sprintf("failed to marshal authorized payload (%s)", err))
	}

	return jsonPayload, nil
}
