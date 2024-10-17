package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	authIface "github.com/circulohealth/sonar-backend/packages/common/authorizer/interfaces"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
)

type GeneratePolicyInput struct {
	Claims   *authorizer.OliveClaims
	Effect   string
	Resource string
}

type HandleRequestInput struct {
	Event  events.APIGatewayCustomAuthorizerRequestTypeRequest
	Jwt    authIface.Jwt
	Logger *zap.Logger
}

func GeneratePolicy(input GeneratePolicyInput) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{}

	if input.Effect != "" && input.Resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   input.Effect,
					Resource: []string{input.Resource},
				},
			},
		}
	}

	// Optional output with custom properties
	authResponse.Context = map[string]interface{}{
		"email": input.Claims.Email,
	}

	return authResponse
}

func handleRequest(input HandleRequestInput) (events.APIGatewayCustomAuthorizerResponse, *exception.SonarError) {
	tokenStr := authorizer.GetAuthorizationToken(input.Event.Headers)

	verifyKey, err := input.Jwt.ParseRSAPublicKeyFromPEM([]byte(authorizer.OlivePublicKey))
	if err != nil {
		input.Logger.Panic(err.Error())
	}

	token, err := input.Jwt.ParseWithClaims(tokenStr, &authorizer.OliveClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		input.Logger.Error("parsing token: " + err.Error())
		return events.APIGatewayCustomAuthorizerResponse{}, exception.NewSonarError(http.StatusBadRequest, "invalid token")
	}

	claims, ok := token.Claims.(*authorizer.OliveClaims)
	if ok && token.Valid {
		input := GeneratePolicyInput{
			Claims:   claims,
			Effect:   "Allow",
			Resource: input.Event.MethodArn,
		}
		return GeneratePolicy(input), nil
	}

	gpi := GeneratePolicyInput{
		Claims:   claims,
		Effect:   "Deny",
		Resource: input.Event.MethodArn,
	}

	policy := GeneratePolicy(gpi)

	return policy, nil
}
