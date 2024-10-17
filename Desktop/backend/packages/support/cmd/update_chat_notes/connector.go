//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"go.uber.org/zap"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("update_chat_notes called")

	id, ok := event.PathParameters["id"]
	if event.PathParameters == nil || !ok {
		config.Logger.Error("parameter path {id} was not found")
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "parameter path {id} was not found")
	}

	sessionId, err := strconv.Atoi(id)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, fmt.Sprintf("sessionId should be an integer (%s)", err))
	}

	var updateChatNotesRequest request.ChatNotesRequest
	err = json.Unmarshal([]byte(config.Event.Body), &updateChatNotesRequest)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger = config.Logger.With(zap.Int("sessionID", sessionId))

	db, err := dao.OpenConnectionWithTablePrefix(dao.Chat)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, fmt.Sprintf("unable to open connection to DB (%s)", err.Error()))
	}

	err = Handler(UpdateChatNotesRequest{
		DB:        db,
		SessionID: sessionId,
		Notes:     updateChatNotesRequest.Notes,
		Logger:    config.Logger,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("update_chat_notes complete")

	return response.OKv2()
}

func main() {
	lambda.Start(connector)
}
