package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	authIface "github.com/circulohealth/sonar-backend/packages/common/authorizer/interfaces"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

type HandleRequestInput struct {
	Context         context.Context
	DB              dynamo.Database
	Event           *events.APIGatewayCustomAuthorizerRequestTypeRequest
	Logger          *zap.Logger
	GroupNamePrefix string
	JwkSet          jwk.Set
	Jwt             authIface.Jwt
}

func handleRequest(input HandleRequestInput) (events.APIGatewayCustomAuthorizerResponse, error) {
	authConfig := &authorizer.AuthorizerConfig{
		Context:         input.Context,
		DB:              input.DB,
		Event:           input.Event,
		Logger:          input.Logger,
		GroupNamePrefix: input.GroupNamePrefix,
		JwkSet:          input.JwkSet,
		Jwt:             input.Jwt,
	}

	return authConfig.BuildResponseForKnownUser()
}
