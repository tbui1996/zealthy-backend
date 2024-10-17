package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dao/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/externalAuth"
)

type HandlerInput struct {
	Event            events.APIGatewayV2HTTPRequest
	Repo             iface.EmailDomainWhitelistRepository
	Idp              cognitoidentityprovideriface.CognitoIdentityProviderAPI
	ClientID         string
	Logger           *zap.Logger
	FamilyName       string
	GivenName        string
	OrganizationName string
}

func handler(input HandlerInput) (*cognitoidentityprovider.SignUpOutput, *exception.SonarError) {
	email, err := externalAuth.GetEmailFromToken(input.Event.Headers)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, err.Error())
	}

	password, err := externalAuth.CreatePassword()
	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "error setting up user for sign up")
	}

	output, err := input.Idp.SignUp(&cognitoidentityprovider.SignUpInput{
		Username:       aws.String(email),
		UserAttributes: []*cognitoidentityprovider.AttributeType{{Name: aws.String("given_name"), Value: aws.String(input.GivenName)}, {Name: aws.String("family_name"), Value: aws.String(input.FamilyName)}, {Name: aws.String("custom:organization"), Value: aws.String(input.OrganizationName)}},
		Password:       aws.String(password),
		ClientId:       aws.String(input.ClientID),
	})

	if err != nil {
		errMsg := handleSignUpErr(err)
		input.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	if output == nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "expected output to not be null")
	}

	return output, nil
}
