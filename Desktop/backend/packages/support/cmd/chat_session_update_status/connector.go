//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"go.uber.org/zap"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	logging.SyncLogger(config.Logger)

	config.Logger.Info("chat_session_update_status called")

	var updateChatSessionRequest request.UpdateChatSessionRequest
	if err := json.Unmarshal([]byte(event.Body), &updateChatSessionRequest); err != nil {
		errMsg := fmt.Sprintf("unable to parse json body, cannot update chat session (%s)", err.Error())
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Chat)
	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to database (%s)", err.Error())
		config.Logger.Error(errMessage)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	config.Logger = config.Logger.With(zap.Bool("sessionState", updateChatSessionRequest.Open))
	config.Logger = config.Logger.With(zap.String("sessionID", updateChatSessionRequest.ID))
	config.Logger.Debug("updating session state")

	err = Handler(UpdateStatusRequest{
		DB:            db,
		SessionId:     updateChatSessionRequest.ID,
		Open:          updateChatSessionRequest.Open,
		RideScheduled: updateChatSessionRequest.RideScheduled,
		Logger:        config.Logger,
	})

	if err != nil {
		errMessage := fmt.Sprintf("error storing item in RDS: (%s)", err)
		config.Logger.Error(errMessage)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	config.Logger.Info("chat_session_update_status completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
