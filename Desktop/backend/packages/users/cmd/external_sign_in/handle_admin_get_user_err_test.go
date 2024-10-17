package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleAdminGetUserErr(t *testing.T) {
	getInput := func(code string) awserr.Error {
		mockErr := new(mocks.Error)
		mockErr.On("Code").Return(code)
		mockErr.On("Error").Return("mock err")
		return mockErr
	}

	tests := []struct {
		inputErr               awserr.Error
		inputCode              string
		expectedSonarException *exception.SonarError
	}{
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeUserNotFoundException),
			inputCode:              cognitoidentityprovider.ErrCodeUserNotFoundException,
			expectedSonarException: exception.NewSonarError(http.StatusNotFound, "mock err"),
		},
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeInvalidParameterException),
			inputCode:              cognitoidentityprovider.ErrCodeInvalidParameterException,
			expectedSonarException: exception.NewSonarError(http.StatusBadRequest, "mock err"),
		},
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeResourceNotFoundException),
			inputCode:              cognitoidentityprovider.ErrCodeResourceNotFoundException,
			expectedSonarException: exception.NewSonarError(http.StatusNotFound, "mock err"),
		},
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeTooManyRequestsException),
			inputCode:              cognitoidentityprovider.ErrCodeTooManyRequestsException,
			expectedSonarException: exception.NewSonarError(http.StatusTooManyRequests, "mock err"),
		},
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeNotAuthorizedException),
			inputCode:              cognitoidentityprovider.ErrCodeNotAuthorizedException,
			expectedSonarException: exception.NewSonarError(http.StatusUnauthorized, "mock err"),
		},
		{
			inputErr:               getInput(cognitoidentityprovider.ErrCodeInternalErrorException),
			inputCode:              cognitoidentityprovider.ErrCodeInternalErrorException,
			expectedSonarException: exception.NewSonarError(http.StatusInternalServerError, "mock err"),
		},
		{
			inputErr:               getInput("unknown code"),
			inputCode:              "unknown code",
			expectedSonarException: exception.NewSonarError(http.StatusInternalServerError, "error when trying to retrieve user (mock err)"),
		},
	}

	for _, test := range tests {
		actualErr := handleAdminGetUserErr(test.inputErr)
		assert.Equal(t, test.expectedSonarException, actualErr)
		test.inputErr.(*mocks.Error).AssertCalled(t, "Code")
		test.inputErr.(*mocks.Error).AssertCalled(t, "Error")
	}
}
