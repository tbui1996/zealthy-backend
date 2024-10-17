package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func handleAdminGetUserErr(err error) *exception.SonarError {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case cognitoidentityprovider.ErrCodeUserNotFoundException:
			// This exception is thrown when a user is not found.
			return exception.NewSonarError(http.StatusNotFound, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidParameterException:
			// This exception is thrown when the Amazon Cognito service encounters an invalid parameter.
			return exception.NewSonarError(http.StatusBadRequest, aerr.Error())
		case cognitoidentityprovider.ErrCodeResourceNotFoundException:
			// This exception is thrown when the Amazon Cognito service cannot find the requested resource.
			return exception.NewSonarError(http.StatusNotFound, aerr.Error())
		case cognitoidentityprovider.ErrCodeTooManyRequestsException:
			// This exception is thrown when the user has made too many requests for a given operation.
			return exception.NewSonarError(http.StatusTooManyRequests, aerr.Error())
		case cognitoidentityprovider.ErrCodeNotAuthorizedException:
			// This exception is thrown when a user is not authorized.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeInternalErrorException:
			// This exception is thrown when Amazon Cognito encounters an internal error.
			return exception.NewSonarError(http.StatusInternalServerError, aerr.Error())
		}
	}

	errMessage := fmt.Sprintf("error when trying to retrieve user (%s)", err)
	return exception.NewSonarError(http.StatusInternalServerError, errMessage)
}
