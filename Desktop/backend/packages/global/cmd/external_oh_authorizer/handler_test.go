package main

import (
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGeneratePolicy(t *testing.T) {
	inputBuilder := func(email, effect, resource string) GeneratePolicyInput {
		return GeneratePolicyInput{
			Claims: &authorizer.OliveClaims{
				Email: email,
			},
			Effect:   effect,
			Resource: resource,
		}
	}

	outputBuilder := func(email, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
		return events.APIGatewayCustomAuthorizerResponse{
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   effect,
						Resource: []string{resource},
					},
				},
			},
			Context: map[string]interface{}{
				"email": email,
			},
		}
	}

	tests := []struct {
		input    GeneratePolicyInput
		expected events.APIGatewayCustomAuthorizerResponse
	}{
		{
			input:    inputBuilder("test@circulohealth.com", "Allow", "MyArn"),
			expected: outputBuilder("test@circulohealth.com", "Allow", "MyArn"),
		},
		{
			input:    inputBuilder("test@circulohealth.com", "Deny", "MyArn"),
			expected: outputBuilder("test@circulohealth.com", "Deny", "MyArn"),
		},
	}

	for _, test := range tests {
		actual := GeneratePolicy(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestHandleRequest(t *testing.T) {
	// setup start
	token := "myToken"
	email := "test@circulohealth.com"

	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: map[string]string{
			"authorization": "Bearer " + token,
		},
		MethodArn: "MyArn",
	}

	mockJwt := mocks.Jwt{}

	mockJwt.On("ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey)).Return(&rsa.PublicKey{}, nil)
	mockJwt.On("ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything).Return(&jwt.Token{
		Claims: &authorizer.OliveClaims{
			Email: email,
		},
		Valid: true,
	}, nil)

	handleRequestInput := HandleRequestInput{
		Event:  event,
		Jwt:    &mockJwt,
		Logger: zap.NewExample(),
	}
	// setup end

	// test assertions for happy path
	expectedResponse := events.APIGatewayCustomAuthorizerResponse{
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{"MyArn"},
				},
			},
		},
		Context: map[string]interface{}{
			"email": "test@circulohealth.com",
		},
	}

	actualResponse, actualErr := handleRequest(handleRequestInput)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Nil(t, actualErr)

	mockJwt.AssertCalled(t, "ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey))
	mockJwt.AssertCalled(t, "ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything)
}

func TestFailGetRSA(t *testing.T) {
	// setup start
	token := "myToken"
	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: map[string]string{
			"authorization": "Bearer " + token,
		},
		MethodArn: "MyArn",
	}

	mockJwt := mocks.Jwt{}

	mockJwt.On("ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey)).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	handleRequestInput := HandleRequestInput{
		Event:  event,
		Jwt:    &mockJwt,
		Logger: zap.NewExample(),
	}

	assert.Panics(t, func() {
		// nolint errcheck
		handleRequest(handleRequestInput)
	}, "It be panicking")

	mockJwt.AssertCalled(t, "ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey))
	mockJwt.AssertNotCalled(t, "ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything)
}

func TestParseWithClaimsFail(t *testing.T) {
	// setup start
	token := "myToken"
	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: map[string]string{
			"authorization": "Bearer " + token,
		},
		MethodArn: "MyArn",
	}

	mockJwt := mocks.Jwt{}

	mockJwt.On("ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey)).Return(&rsa.PublicKey{}, nil)
	mockJwt.On("ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	handleRequestInput := HandleRequestInput{
		Event:  event,
		Jwt:    &mockJwt,
		Logger: zap.NewExample(),
	}
	// setup end

	_, actualErr := handleRequest(handleRequestInput)
	assert.NotNil(t, actualErr)

	mockJwt.AssertCalled(t, "ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey))
	mockJwt.AssertCalled(t, "ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything)
}

func TestHandleRequestDeny(t *testing.T) {
	// setup start
	token := "myToken"
	email := "test@circulohealth.com"

	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: map[string]string{
			"authorization": "Bearer " + token,
		},
		MethodArn: "MyArn",
	}

	mockJwt := mocks.Jwt{}

	mockJwt.On("ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey)).Return(&rsa.PublicKey{}, nil)
	mockJwt.On("ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything).Return(&jwt.Token{
		Claims: &authorizer.OliveClaims{
			Email: email,
		},
		Valid: false,
	}, nil)

	handleRequestInput := HandleRequestInput{
		Event:  event,
		Jwt:    &mockJwt,
		Logger: zap.NewExample(),
	}
	// setup end

	// test assertions for happy path
	expectedResponse := events.APIGatewayCustomAuthorizerResponse{
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Deny",
					Resource: []string{"MyArn"},
				},
			},
		},
		Context: map[string]interface{}{
			"email": "test@circulohealth.com",
		},
	}

	actualResponse, actualErr := handleRequest(handleRequestInput)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Nil(t, actualErr)

	mockJwt.AssertCalled(t, "ParseRSAPublicKeyFromPEM", []byte(authorizer.OlivePublicKey))
	mockJwt.AssertCalled(t, "ParseWithClaims", token, &authorizer.OliveClaims{}, mock.Anything)
}
