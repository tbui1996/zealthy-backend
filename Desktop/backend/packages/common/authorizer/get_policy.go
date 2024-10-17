package authorizer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"go.uber.org/zap"
)

type GetPolicyInput struct {
	Claims          *AccessTokenClaims
	DB              dynamo.Database
	GroupNamePrefix string
	Logger          *zap.Logger
}

func GetPolicy(input GetPolicyInput) (*events.APIGatewayCustomAuthorizerPolicy, error) {
	// Get first group from claim or error if none
	// grab the policy from dynamodb
	// put into a policy and return

	if input.Claims.CognitoGroups == nil || len(input.Claims.CognitoGroups) == 0 {
		errMessage := "expected cognito groups to be passed in access token"
		input.Logger.Error(errMessage)
		return nil, fmt.Errorf("expected cognito groups to be passed in access token")
	}

	var group string
	for _, cognitoGroup := range input.Claims.CognitoGroups {
		if strings.Contains(cognitoGroup, input.GroupNamePrefix) {
			group = strings.TrimPrefix(cognitoGroup, fmt.Sprintf("%s_", input.GroupNamePrefix))
		}
	}

	if group == "" {
		return nil, fmt.Errorf("expected to find a group in cognito groups for prefix (%s)", input.GroupNamePrefix)
	}

	output, err := input.DB.Get(&dynamodb.GetItemInput{
		TableName: aws.String(dynamo.SonarGroupPolicy),
		Key: map[string]*dynamodb.AttributeValue{
			"group": {
				// grab the first group (expecting only 1 group associated to the user)
				S: aws.String(group),
			},
		},
	})
	if err != nil {
		input.Logger.Error("failed to get policy from dynamodb")
		return nil, err
	}

	policyJSONAttribute, ok := output.Item["policy"]
	if !ok {
		return nil, fmt.Errorf("expected policy to exist in dynamodb table for group %s", input.Claims.CognitoGroups[0])
	}

	if policyJSONAttribute.S == nil {
		return nil, fmt.Errorf("policy document exists but could not dereference value for group %s", input.Claims.CognitoGroups[0])
	}

	var policy events.APIGatewayCustomAuthorizerPolicy
	if err = json.Unmarshal([]byte(*policyJSONAttribute.S), &policy); err != nil {
		input.Logger.Error("failed to unmarshall policy: (" + *policyJSONAttribute.S + "): " + err.Error())
		return nil, err
	}

	return &policy, nil
}
