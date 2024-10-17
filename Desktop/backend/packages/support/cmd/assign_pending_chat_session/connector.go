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

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}
	logging.SyncLogger(config.Logger)

	var assignPendingChatSessionRequest request.AssignPendingChatSessionRequestInternal
	err = json.Unmarshal([]byte(config.Event.Body), &assignPendingChatSessionRequest)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	repo, err := chatSession.NewRDSChatSessionRepositoryWithSession(config.Session)
	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to database (%s)", err.Error())
		config.Logger.Error(errMessage)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	config.Logger.Info("assign_pending_chat_session called")

	body, sErr := Handler(config.Logger, repo, assignPendingChatSessionRequest)
	if sErr != nil {
		config.Logger.Error("handling assign_pending_chat_session request: " + sErr.Error())
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("assign_pending_chat_session completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
