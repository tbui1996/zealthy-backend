package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer/mocks"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	cMocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHandleRequest(t *testing.T) {
	// setup start
	groupName := "program_manager"
	groupNamePrefix := "internals"
	expectedResponse := events.APIGatewayCustomAuthorizerResponse{
		Context: map[string]interface{}{
			"userID": "testuser",
			"group":  "internals_program_manager",
		},
	}

	mockDB := new(dynamo.MockDatabase)
	mockJwkSet := new(cMocks.Set)
	mockJwt := new(mocks.Jwt)

	mockJwt.On("ParseWithClaims", "test", mock.Anything, mock.Anything).Return(&jwt.Token{
		Claims: &authorizer.AccessTokenClaims{
			Username:      "testuser",
			CognitoGroups: []string{"internals_program_manager"},
		},
		Valid: true,
	}, nil /** error */)

	itemIn := &dynamodb.GetItemInput{
		TableName: aws.String(dynamo.SonarGroupPolicy),
		Key: map[string]*dynamodb.AttributeValue{
			"group": {
				S: aws.String(groupName),
			},
		},
	}

	itemOut := &dynamodb.GetItemOutput{
		ConsumedCapacity: &dynamodb.ConsumedCapacity{},
		Item: map[string]*dynamodb.AttributeValue{
			"policy": {
				S: aws.String("{}"),
			},
		},
	}

	mockDB.On("Get", itemIn).Return(itemOut, nil)

	input := HandleRequestInput{
		Context: context.TODO(),
		DB:      mockDB,
		Event: &events.APIGatewayCustomAuthorizerRequestTypeRequest{
			Headers: map[string]string{
				"authorization": "test",
			},
		},
		Logger:          zap.NewExample(),
		GroupNamePrefix: groupNamePrefix,
		JwkSet:          mockJwkSet,
		Jwt:             mockJwt,
	}
	// setup end

	// test assertions for happy path
	actualResponse, actualErr := handleRequest(input)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, nil, actualErr)

	assert.IsType(t, new(dynamo.MockDatabase), input.DB)
	mockDB, _ = input.DB.(*dynamo.MockDatabase)
	mockDB.AssertCalled(t, "Get", itemIn)

	assert.IsType(t, new(mocks.Jwt), input.Jwt)
	mockJwt, _ = input.Jwt.(*mocks.Jwt)
	mockJwt.AssertCalled(t, "ParseWithClaims", "test", mock.Anything, mock.Anything)
}
