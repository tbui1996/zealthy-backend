package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	cMocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
	input HandlerInput
}

func (suite *HandlerSuite) SetupTest() {
	suite.input = HandlerInput{
		Logger:   zaptest.NewLogger(suite.T()),
		Idp:      new(cMocks.CognitoIdentityProviderAPI),
		Email:    "email@test.com",
		ClientID: "client_id",
		PoolID:   "pool_id",
	}
}

func (suite *HandlerSuite) TestErrGettingUser() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(nil, errors.New("idpError"))

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "error when trying to retrieve user (idpError)"), err)
	suite.Equal("", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
}

func (suite *HandlerSuite) TestIdpReturnsNilUser() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(nil, nil)

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "expected user to not be null"), err)
	suite.Equal("", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
}

func (suite *HandlerSuite) TestIdpReturnsUserNotEnabled() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(&cognitoidentityprovider.AdminGetUserOutput{
		Enabled: aws.Bool(false),
	}, nil)

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusUnauthorized, "user is registered, awaiting confirmation and role assignment"), err)
	suite.Equal("", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
}

func (suite *HandlerSuite) TestIdpReturnsUserEnabledButNotConfirmed() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(&cognitoidentityprovider.AdminGetUserOutput{
		Enabled:    aws.Bool(true),
		UserStatus: aws.String(cognitoidentityprovider.UserStatusTypeUnconfirmed),
	}, nil)

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "unknown user status: UNCONFIRMED"), err)
	suite.Equal("", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
}

func (suite *HandlerSuite) TestIdpErrInitiatingAuth() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(&cognitoidentityprovider.AdminGetUserOutput{
		Enabled:    aws.Bool(true),
		UserStatus: aws.String(cognitoidentityprovider.UserStatusTypeConfirmed),
	}, nil)

	initAuthInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(suite.input.Email),
		},
		ClientId:   aws.String(suite.input.ClientID),
		UserPoolId: aws.String(suite.input.PoolID),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminInitiateAuth", initAuthInput).Return(nil, errors.New("idpError"))

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "error initiating auth (idpError)"), err)
	suite.Equal("", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminInitiateAuth", initAuthInput)
}

func (suite *HandlerSuite) TestHappyPath() {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(suite.input.PoolID),
		Username:   aws.String(suite.input.Email),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminGetUser", getUserInput).Return(&cognitoidentityprovider.AdminGetUserOutput{
		Enabled:    aws.Bool(true),
		UserStatus: aws.String(cognitoidentityprovider.UserStatusTypeConfirmed),
	}, nil)

	initAuthInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(suite.input.Email),
		},
		ClientId:   aws.String(suite.input.ClientID),
		UserPoolId: aws.String(suite.input.PoolID),
	}

	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("AdminInitiateAuth", initAuthInput).Return(&cognitoidentityprovider.AdminInitiateAuthOutput{
		AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
			IdToken:      aws.String("id_token"),
			AccessToken:  aws.String("access_token"),
			RefreshToken: aws.String("refresh_token"),
		},
	}, nil)

	output, err := handler(suite.input)
	suite.Nil(err)
	suite.Equal("{\"id\":\"id_token\",\"token\":\"access_token\",\"refreshToken\":\"refresh_token\"}", output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminGetUser", getUserInput)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminInitiateAuth", initAuthInput)
}

func TestExternalSignUp(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
