package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandlerSuite struct {
	suite.Suite
	input HandlerInput
}

func (suite *HandlerSuite) SetupTest() {
	suite.input = HandlerInput{
		UserName:     "test_user",
		PoolID:       "pool_id",
		Idp:          new(mocks.CognitoIdentityProviderAPI),
		Logger:       zaptest.NewLogger(suite.T()),
		DefaultGroup: "externals_guest",
	}
}

func (suite *HandlerSuite) TestErrAdminAddUserToGroup() {
	suite.input.Idp.(*mocks.CognitoIdentityProviderAPI).On("AdminAddUserToGroup", &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String("externals_guest"),
		Username:   aws.String(suite.input.UserName),
		UserPoolId: aws.String(suite.input.PoolID),
	}).Return(nil, errors.New("idp err"))

	actualErr := handler(suite.input)
	suite.Equal(http.StatusBadRequest, actualErr.StatusCode)

	suite.input.Idp.(*mocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminAddUserToGroup", &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String("externals_guest"),
		Username:   aws.String(suite.input.UserName),
		UserPoolId: aws.String(suite.input.PoolID),
	})
}

func (suite *HandlerSuite) TestHappyPath() {
	suite.input.Idp.(*mocks.CognitoIdentityProviderAPI).On("AdminAddUserToGroup", &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String("externals_guest"),
		Username:   aws.String(suite.input.UserName),
		UserPoolId: aws.String(suite.input.PoolID),
	}).Return(nil, nil)

	actualErr := handler(suite.input)

	suite.Nil(actualErr)
	suite.input.Idp.(*mocks.CognitoIdentityProviderAPI).AssertCalled(suite.T(), "AdminAddUserToGroup", &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String("externals_guest"),
		Username:   aws.String(suite.input.UserName),
		UserPoolId: aws.String(suite.input.PoolID),
	})
}

func TestPostConfirmation(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
