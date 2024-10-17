//go:build !test
// +build !test

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	logging.SyncLogger(config.Logger)

	config.Logger.Info("chat_session_get called")

	repo, err := chatSession.NewRDSChatSessionRepositoryWithSession(config.Session)

	userID, ok := config.Event.RequestContext.Authorizer.Lambda["userID"].(string)
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "expected a valid userID from the authorizer")
	}

	group, ok := config.Event.RequestContext.Authorizer.Lambda["group"].(string)
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "expected a valid group from the authorizer")
	}

	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to database (%s)", err.Error())
		config.Logger.Error(errMessage)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	config.Logger.Debug("get chat sessions")

	body, err := Handler(userID, group, config.Logger, repo)

	if err != nil {
		config.Logger.Error(err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("chat_session_get completed, returning response")

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
