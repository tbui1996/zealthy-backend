package main

import (
	"errors"
	"net/http"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	cMocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dao/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/externalAuth"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	input            HandlerInput
	idpInput         *cognitoidentityprovider.SignUpInput
	password         string
	email            string
	givenName        string
	familyName       string
	organizationName string
}

func (suite *HandlerSuite) SetupTest() {
	suite.password = "fake_pa$$word"
	suite.email = "test@email.com"
	suite.idpInput = &cognitoidentityprovider.SignUpInput{
		Username:       aws.String(suite.email),
		Password:       aws.String(suite.password),
		ClientId:       aws.String("clientid"),
		UserAttributes: []*cognitoidentityprovider.AttributeType{{Name: aws.String("given_name"), Value: aws.String(suite.givenName)}, {Name: aws.String("family_name"), Value: aws.String(suite.familyName)}, {Name: aws.String("custom:organization"), Value: aws.String(suite.organizationName)}},
	}

	suite.input = HandlerInput{
		Event:    events.APIGatewayV2HTTPRequest{},
		Repo:     new(mocks.EmailDomainWhitelistRepository),
		Idp:      new(cMocks.CognitoIdentityProviderAPI),
		ClientID: "clientid",
		Logger:   zaptest.NewLogger(suite.T()),
	}
}

func (suite *HandlerSuite) TestErrGettingEmailFromHeaders() {
	externalAuth.GetEmailFromToken = func(map[string]string) (string, error) {
		return "", errors.New("error here")
	}

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusBadRequest, "error here"), err)
	suite.Nil(output)
}

func (suite *HandlerSuite) TestErrorCreatingPassword() {
	externalAuth.GetEmailFromToken = func(map[string]string) (string, error) {
		return suite.email, nil
	}

	externalAuth.CreatePassword = func() (string, error) {
		return "", exception.NewSonarError(http.StatusInternalServerError, "error")
	}

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "error setting up user for sign up"), err)
	suite.Nil(output)
}

func (suite *HandlerSuite) TestErrorSigningUp() {
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("SignUp", suite.idpInput).Return(nil, errors.New("idpError"))

	externalAuth.GetEmailFromToken = func(map[string]string) (string, error) {
		return suite.email, nil
	}

	externalAuth.CreatePassword = func() (string, error) {
		return suite.password, nil
	}

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "error when trying to sign up user (idpError)"), err)
	suite.Nil(output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "SignUp", suite.idpInput)
}

func (suite *HandlerSuite) TestErrorSigningUpOutputNil() {
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("SignUp", suite.idpInput).Return(nil, nil)

	externalAuth.GetEmailFromToken = func(map[string]string) (string, error) {
		return suite.email, nil
	}

	externalAuth.CreatePassword = func() (string, error) {
		return suite.password, nil
	}

	output, err := handler(suite.input)
	suite.Equal(exception.NewSonarError(http.StatusInternalServerError, "expected output to not be null"), err)
	suite.Nil(output)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "SignUp", suite.idpInput)
}

func (suite *HandlerSuite) TestHappyPath() {
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).On("SignUp", suite.idpInput).Return(&cognitoidentityprovider.SignUpOutput{}, nil)

	externalAuth.GetEmailFromToken = func(map[string]string) (string, error) {
		return suite.email, nil
	}

	externalAuth.CreatePassword = func() (string, error) {
		return suite.password, nil
	}

	output, err := handler(suite.input)
	suite.Equal(&cognitoidentityprovider.SignUpOutput{}, output)
	suite.Nil(err)
	suite.input.Idp.(*cMocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "SignUp", suite.idpInput)
}

func TestExternalSignUp(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
