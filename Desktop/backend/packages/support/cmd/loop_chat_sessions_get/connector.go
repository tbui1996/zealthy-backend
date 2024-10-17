//go:build !test
// +build !test

package main

import (
	"context"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	logging.SyncLogger(config.Logger)

	if config.Event.PathParameters == nil {
		return exception.NewSonarError(http.StatusBadRequest, "expected path parameters to exist.").ToAPIGatewayV2HTTPResponse(), nil
	}

	userID, ok := config.Event.PathParameters["id"]
	if !ok {
		return exception.NewSonarError(http.StatusBadRequest, "expected path parameter id (userID) to exist").ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Debug("getting entities by external id")
	chatSessionRepo, err := chatSession.NewRDSChatSessionRepositoryWithSession(config.Session)

	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to db (%s)", err)
		config.Logger.Error(errMessage)
		return exception.NewSonarError(http.StatusInternalServerError, errMessage).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("loop_chat_session_get called")

	body, sErr := Handler(config.Logger, chatSessionRepo, userID)
	if sErr != nil {
		config.Logger.Error("handling loop_chat_session_get request: " + sErr.Error())
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("loop_chat_session_get completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
