package main

import (
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
)

func Handler(
	config *requestConfig.APIGatewayWebsocketProxyRequest,
	repo iface.ChatSessionRepository,
	client *router.Session,
) (*events.APIGatewayProxyResponse, *exception.SonarError) {
	var supportRequest request.SupportRequestSend
	err := json.Unmarshal([]byte(config.Event.Body), &supportRequest)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, fmt.Sprintf("unable to unmarshal body: %+v", config.Event.Body))
	}

	switch supportRequest.Payload.Type {
	case "chat":
		config.Logger.Debug("handling chat message")
		err = HandleChat(config, supportRequest.Payload.Message, repo, client)

		if err != nil {
			config.Logger.Error("error: " + err.Error())
			return nil, exception.NewSonarError(http.StatusInternalServerError, err.Error())
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "OK",
		}, nil

	case "read_receipt":
		config.Logger.Debug("handling read receipt message")
		err = chatHelper.HandleReadReceipt(supportRequest.Payload.Message, time.Now().Unix(), repo)

		if err != nil {
			return nil, exception.NewSonarError(http.StatusInternalServerError, err.Error())
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "OK",
		}, nil

	case "typing":
		config.Logger.Debug("handling typing event")
		err = HandleTyping(config, supportRequest.Payload.Message, repo, client)

		if err != nil {
			return nil, exception.NewSonarError(http.StatusBadRequest, err.Error())
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "OK",
		}, nil

	default:
		return nil, exception.NewSonarError(http.StatusBadRequest, fmt.Sprintf("unsupported request type %s", supportRequest.Payload.Type))
	}
}
