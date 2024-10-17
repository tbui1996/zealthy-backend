package authorizer

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type DBIO struct {
	Input  *dynamodb.GetItemInput
	Output *dynamodb.GetItemOutput
}

func getPolicyInputBuilder(groupPrefix, groupName string, claims *AccessTokenClaims, dbErr error) (GetPolicyInput, DBIO) {
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

	return GetPolicyInput{
			Claims:          claims,
			DB:              mockDB,
			GroupNamePrefix: groupPrefix,
			Logger:          zap.NewExample(),
		}, DBIO{
			Input:  itemIn,
			Output: itemOut,
		}
}

func TestGetPolicy(t *testing.T) {
	tests := []struct {
		inputClaims          *AccessTokenClaims
		inputGroupName       string
		inputGroupNamePrefix string
		outputErr            error
		expectedPolicy       *events.APIGatewayCustomAuthorizerPolicy
	}{
		// happy path internals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"not_a_real_group", "internals_program_manager"}},
			inputGroupName:       "program_manager",
			inputGroupNamePrefix: "internals",
			outputErr:            nil,
			expectedPolicy:       &events.APIGatewayCustomAuthorizerPolicy{},
		},
		// internal trying to get policy for externals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"internals_supervisor"}},
			inputGroupName:       "supervisor",
			inputGroupNamePrefix: "externals",
			outputErr:            errors.New("expected to find a group in cognito groups for prefix (externals)"),
			expectedPolicy:       nil,
		},
		// happy path externals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"not_a_real_group", "externals_supervisor"}},
			inputGroupName:       "supervisor",
			inputGroupNamePrefix: "externals",
			outputErr:            nil,
			expectedPolicy:       &events.APIGatewayCustomAuthorizerPolicy{},
		},
		// external trying to get policy for internals
		{
			inputClaims:          &AccessTokenClaims{CognitoGroups: []string{"externals_supervisor"}},
			inputGroupName:       "program_manager",
			inputGroupNamePrefix: "internals",
			outputErr:            errors.New("expected to find a group in cognito groups for prefix (internals)"),
			expectedPolicy:       nil,
		},
	}

	for _, test := range tests {
		getPolicyInput, dbio := getPolicyInputBuilder(test.inputGroupNamePrefix, test.inputGroupName, test.inputClaims, test.outputErr)

		actualPolicy, actualErr := GetPolicy(getPolicyInput)
		assert.Equal(t, test.expectedPolicy, actualPolicy)
		assert.Equal(t, test.outputErr, actualErr)

		assert.IsType(t, new(dynamo.MockDatabase), getPolicyInput.DB)

		mockDB, _ := getPolicyInput.DB.(*dynamo.MockDatabase)

		if test.outputErr == nil {
			mockDB.AssertCalled(t, "Get", dbio.Input)
		}
	}
}
