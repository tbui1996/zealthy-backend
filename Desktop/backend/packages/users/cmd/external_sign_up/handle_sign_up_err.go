package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func handleSignUpErr(err error) *exception.SonarError {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case cognitoidentityprovider.ErrCodeUsernameExistsException:
			return exception.NewSonarError(http.StatusConflict, aerr.Error())
		case cognitoidentityprovider.ErrCodeResourceNotFoundException:
			// This exception is thrown when the Amazon Cognito service cannot find the requested resource.
			return exception.NewSonarError(http.StatusNotFound, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidParameterException:
			// This exception is thrown when the Amazon Cognito service encounters an invalid parameter.
			return exception.NewSonarError(http.StatusBadRequest, aerr.Error())
		case cognitoidentityprovider.ErrCodeUnexpectedLambdaException:
			// This exception is thrown when the Amazon Cognito service encounters an unexpected exception with the Lambda service.
			return exception.NewSonarError(http.StatusInternalServerError, aerr.Error())
		case cognitoidentityprovider.ErrCodeUserLambdaValidationException:
			// This exception is thrown when the Amazon Cognito service encounters a user validation exception with the Lambda service.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeNotAuthorizedException:
			// This exception is thrown when a user is not authorized.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidPasswordException:
			// This exception is thrown when the Amazon Cognito service encounters an invalid password.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidLambdaResponseException:
			// This exception is thrown when the Amazon Cognito service encounters an invalid Lambda response.
			return exception.NewSonarError(http.StatusInternalServerError, aerr.Error())
		case cognitoidentityprovider.ErrCodeTooManyRequestsException:
			// This exception is thrown when the user has made too many requests for a given operation.
			return exception.NewSonarError(http.StatusTooManyRequests, aerr.Error())
		case cognitoidentityprovider.ErrCodeInternalErrorException:
			// This exception is thrown when Amazon Cognito encounters an internal error.
			return exception.NewSonarError(http.StatusInternalServerError, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidSmsRoleAccessPolicyException:
			// This exception is returned when the role provided for SMS configuration does not have permission to publish using Amazon SNS.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidSmsRoleTrustRelationshipException:
			// This exception is thrown when the trust relationship is invalid for the role provided for SMS configuration. This can happen if you do not trust cognito-idp.amazonaws.com or the external ID provided in the role does not match what is provided in the SMS configuration for the user pool.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeInvalidEmailRoleAccessPolicyException:
			// This exception is thrown when Amazon Cognito is not allowed to use your email identity. HTTP status code: 400.
			return exception.NewSonarError(http.StatusUnauthorized, aerr.Error())
		case cognitoidentityprovider.ErrCodeCodeDeliveryFailureException:
			// This exception is thrown when a verification code fails to deliver successfully.
			return exception.NewSonarError(http.StatusInternalServerError, aerr.Error())
		}
	}

	errMessage := fmt.Sprintf("error when trying to sign up user (%s)", err)
	return exception.NewSonarError(http.StatusInternalServerError, errMessage)
}
