//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}
	logging.SyncLogger(config.Logger)

	config.Logger.Info("chat_session_create called")

	var createChatSessionRequest request.ChatSessionRequestInternal
	if err := json.Unmarshal([]byte(config.Event.Body), &createChatSessionRequest); err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "unable to parse json body, cannot create form"+err.Error())
	}

	repo, err := chatSession.NewRDSChatSessionRepositoryWithSession(config.Session)

	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to database (%s)", err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	bodyBytes, sErr := Handler(config.Logger, repo, createChatSessionRequest)
	if sErr != nil {
		config.Logger.Error("chat_session_create failed: " + sErr.Error())
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("chat_session_create completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
