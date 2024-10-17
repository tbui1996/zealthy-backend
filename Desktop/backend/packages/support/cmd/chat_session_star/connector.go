//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("chat_session_star called")

	var chatSessionStarRequest request.ChatSessionStarRequest
	if mErr := json.Unmarshal([]byte(config.Event.Body), &chatSessionStarRequest); mErr != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, mErr.Error())
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Chat)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, fmt.Sprintf("unable to open connection to DB (%s)", err.Error()))
	}

	err = Handler(SubmitChatSessionStarRequest{
		DB:        db,
		SessionID: chatSessionStarRequest.SessionID,
		OnStar:    chatSessionStarRequest.OnStar,
		Logger:    config.Logger,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("chat_session_star complete")

	return response.OKv2()
}

func main() {
	lambda.Start(connector)
}
