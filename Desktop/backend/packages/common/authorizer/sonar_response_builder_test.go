package authorizer

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer/mocks"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	cMocks "github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestBuildResponseForKnownUser(t *testing.T) {
	// setup start
	groupName := "program_manager"
	groupNamePrefix := "internals"

	mockDB := new(dynamo.MockDatabase)
	mockJwkSet := new(cMocks.Set)
	mockJwt := new(mocks.Jwt)

	mockJwt.On("ParseWithClaims", "test", mock.Anything, mock.Anything).Return(&jwt.Token{
		Claims: &AccessTokenClaims{
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

	ac := AuthorizerConfig{Context: context.TODO(), DB: mockDB, Event: &events.APIGatewayCustomAuthorizerRequestTypeRequest{Headers: map[string]string{"authorization": "test"}}, Logger: zap.NewExample(), GroupNamePrefix: groupNamePrefix, JwkSet: mockJwkSet, Jwt: mockJwt}
	// setup end

	// test assertions for happy path
	expectedResponse := events.APIGatewayCustomAuthorizerResponse{
		Context: map[string]interface{}{
			"userID": "testuser",
			"group":  "internals_program_manager",
		},
	}

	actualResponse, actualErr := ac.BuildResponseForKnownUser()
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, nil, actualErr)

	assert.IsType(t, new(dynamo.MockDatabase), ac.DB)
	mockDB, _ = ac.DB.(*dynamo.MockDatabase)
	mockDB.AssertCalled(t, "Get", itemIn)

	assert.IsType(t, new(mocks.Jwt), ac.Jwt)
	mockJwt, _ = ac.Jwt.(*mocks.Jwt)
	mockJwt.AssertCalled(t, "ParseWithClaims", "test", mock.Anything, mock.Anything)
}

func buildResponseBuilder(groupPrefix, groupName string, claims *AccessTokenClaims, dbErr error) (BuildResponseInput, DBIO) {
	mockDB := new(dynamo.MockDatabase)

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

	mockDB.On("Get", itemIn).Return(itemOut, dbErr)

	return BuildResponseInput{
			Claims:          claims,
			DB:              mockDB,
			GroupNamePrefix: groupPrefix,
			Logger:          zap.NewExample(),
		}, DBIO{
			Input:  itemIn,
			Output: itemOut,
		}
}

func TestBuildResponse(t *testing.T) {
	tests := []struct {
		inputClaims          *AccessTokenClaims
		inputGroupName       string
		inputGroupNamePrefix string
		dbErr                error
		expectedResponse     *events.APIGatewayCustomAuthorizerResponse
	}{
		// happy path internals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"not_a_real_group", "internals_program_manager"}, Username: "testuser"},
			inputGroupName:       "program_manager",
			inputGroupNamePrefix: "internals",
			dbErr:                nil,
			expectedResponse:     &events.APIGatewayCustomAuthorizerResponse{Context: map[string]interface{}{"userID": "testuser", "group": "internals_program_manager"}},
		},
		// internal trying to get policy for externals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"internals_supervisor"}, Username: "testuser"},
			inputGroupName:       "supervisor",
			inputGroupNamePrefix: "externals",
			dbErr:                errors.New("expected to find a group in cognito groups for prefix (externals)"),
			expectedResponse:     nil,
		},
		// happy path externals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"not_a_real_group", "externals_supervisor"}, Username: "testuser"},
			inputGroupName:       "supervisor",
			inputGroupNamePrefix: "externals",
			dbErr:                nil,
			expectedResponse:     &events.APIGatewayCustomAuthorizerResponse{Context: map[string]interface{}{"userID": "testuser", "group": "externals_supervisor"}},
		},
		// external trying to get policy for internals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"externals_supervisor"}},
			inputGroupName:       "program_manager",
			inputGroupNamePrefix: "internals",
			dbErr:                errors.New("expected to find a group in cognito groups for prefix (internals)"),
			expectedResponse:     nil,
		},
	}

	for _, test := range tests {
		getPolicyInput, dbio := buildResponseBuilder(test.inputGroupNamePrefix, test.inputGroupName, test.inputClaims, test.dbErr)

		actualResponse, actualErr := BuildResponse(getPolicyInput)
		assert.Equal(t, test.expectedResponse, actualResponse)
		assert.Equal(t, test.dbErr, actualErr)

		assert.IsType(t, new(dynamo.MockDatabase), getPolicyInput.DB)

		mockDB, _ := getPolicyInput.DB.(*dynamo.MockDatabase)

		if test.dbErr == nil {
			mockDB.AssertCalled(t, "Get", dbio.Input)
		}
	}
}
