package authorizer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer/interfaces"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

type AuthorizerConfig struct {
	Context         context.Context
	DB              dynamo.Database
	Event           *events.APIGatewayCustomAuthorizerRequestTypeRequest
	Logger          *zap.Logger
	GroupNamePrefix string
	JwkSet          jwk.Set
	Jwt             interfaces.Jwt
}

func (ac *AuthorizerConfig) BuildResponseForKnownUser() (events.APIGatewayCustomAuthorizerResponse, error) {
	token := GetToken(*ac.Event)
	if token == "" {
		ac.Logger.Error("could not find token in expected locations")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("could not find token in expected locations")
	}

	jwt, err := ac.Jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		key, keyExists := ac.JwkSet.LookupKeyID(kid)
		if !keyExists {
			return nil, fmt.Errorf("key %v not found in JWKS", kid)
		}

		var raw interface{}
		return raw, key.Raw(&raw)
	})
	if err != nil {
		ac.Logger.Error("parsing sonar token: " + err.Error())
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("error: invalid token")
	}

	if claims, ok := jwt.Claims.(*AccessTokenClaims); ok && jwt.Valid {
		buildResponseInput := BuildResponseInput{
			Claims:          claims,
			GroupNamePrefix: ac.GroupNamePrefix,
			DB:              ac.DB,
			Logger:          ac.Logger,
		}

		response, err := BuildResponse(buildResponseInput)

		if err != nil {
			ac.Logger.Error("failed to build response: " + err.Error())
			return events.APIGatewayCustomAuthorizerResponse{}, err
		}

		if response == nil {
			return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("expected response to exist")
		}

		return *response, nil
	}

	ac.Logger.Error("invalid token")
	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("error: invalid token")
}

type BuildResponseInput struct {
	Claims          *AccessTokenClaims
	GroupNamePrefix string
	DB              dynamo.Database
	Logger          *zap.Logger
}

func BuildResponse(input BuildResponseInput) (*events.APIGatewayCustomAuthorizerResponse, error) {
	authResponse := events.APIGatewayCustomAuthorizerResponse{}

	getPolicyInput := GetPolicyInput{
		Claims:          input.Claims,
		DB:              input.DB,
		GroupNamePrefix: input.GroupNamePrefix,
		Logger:          input.Logger,
	}
	policy, err := GetPolicy(getPolicyInput)

	if err != nil {
		input.Logger.Error("failed to get policy: " + err.Error())
		return nil, err
	}

	if policy == nil {
		return nil, errors.New("expected policy reference to exist")
	}

	input.Logger.Debug("received policy")
	authResponse.PolicyDocument = *policy

	var group string
	for index, value := range input.Claims.CognitoGroups {
		if strings.HasPrefix(value, fmt.Sprintf("%s_", input.GroupNamePrefix)) {
			group = input.Claims.CognitoGroups[index]
		}
	}

	// Optional output with custom properties
	authResponse.Context = map[string]interface{}{
		"userID": input.Claims.Username,
		"group":  group,
	}

	return &authResponse, nil
}
