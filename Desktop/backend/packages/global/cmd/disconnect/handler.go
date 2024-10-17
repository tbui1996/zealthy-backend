package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
	"go.uber.org/zap"
)

type HandleRequestInput struct {
	DB     dynamo.Database
	Event  events.APIGatewayWebsocketProxyRequest
	Logger *zap.Logger
}

func handleRequest(input HandleRequestInput) *exception.SonarError {
	context, sErr := validate.ConnectionEvent(input.Event)
	if sErr != nil {
		return sErr
	}

	input.Logger = input.Logger.With(zap.String("connectionID", context.ConnectionID))
	input.Logger.Debug("disconnecting client")

	sErr = removeConnectionInfo(RemoveConnectionInfoInput{
		Context: *context,
		DB:      input.DB,
	})
	if sErr != nil {
		return sErr
	}

	return nil
}
