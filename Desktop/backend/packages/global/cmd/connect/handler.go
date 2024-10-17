package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/events/iface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
	"go.uber.org/zap"
	"net/http"
)

type HandleRequestInput struct {
	DB             dynamo.Database
	Event          events.APIGatewayWebsocketProxyRequest
	Logger         *zap.Logger
	EventPublisher iface.EventPublisher
}

func handleRequest(input HandleRequestInput) *exception.SonarError {
	context, sErr := validate.ConnectionEvent(input.Event)
	if sErr != nil {
		return sErr
	}

	input.Logger = input.Logger.With(zap.String("connectionID", context.ConnectionID))
	input.Logger.Debug("connecting client")

	sErr = storeConnectionInfo(StoreConnectionInfoInput{
		Context: *context,
		DB:      input.DB,
	})

	if sErr != nil {
		return sErr
	}

	err := input.EventPublisher.PublishConnectionCreatedEvent(context.UserID, eventconstants.GLOBAL_SERVICE)

	if err != nil {
		input.Logger.Error(fmt.Sprintf("unable to publish event: %s", err.Error()))
		return exception.NewSonarError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
