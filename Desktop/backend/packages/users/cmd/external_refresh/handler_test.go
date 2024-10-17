package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	inputBuilder := func(idpOut *cognitoidentityprovider.InitiateAuthOutput, idpErr error) HandlerInput {
		mockIdp := new(mocks.CognitoIdentityProviderAPI)
		mockIdp.On("InitiateAuth", mock.Anything).Return(idpOut, idpErr)

		return HandlerInput{
			ClientID:     "clientID",
			Idp:          mockIdp,
			RefreshToken: "refresh_token",
		}
	}

	tests := []struct {
		input           HandlerInput
		expectedPayload []byte
		expectedErr     *exception.SonarError
	}{
		{
			// valid and successful request
			input: inputBuilder(&cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
					AccessToken: aws.String("access_token"),
					IdToken:     aws.String("id_token"),
				},
			}, nil),
			expectedPayload: []byte{0x7b, 0x22, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x69, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2c, 0x22, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3a, 0x22, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x7d},
			expectedErr:     nil,
		},
		{
			// InitiateAuth errors
			input:           inputBuilder(nil, errors.New("test")),
			expectedPayload: nil,
			expectedErr:     exception.NewSonarError(http.StatusInternalServerError, "unhandled error when trying to retrieve user list (test)"),
		},
		{
			// InitiateAuth returned output but ValidateInitiateAuthOutput fails
			input:           inputBuilder(&cognitoidentityprovider.InitiateAuthOutput{}, nil),
			expectedPayload: nil,
			expectedErr:     exception.NewSonarError(http.StatusBadRequest, "expected authentication result to not be null"),
		},
	}

	for _, test := range tests {
		actualPayload, actualErr := handler(test.input)
		assert.Equal(t, test.expectedPayload, actualPayload)
		assert.Equal(t, test.expectedErr, actualErr)
	}
}
