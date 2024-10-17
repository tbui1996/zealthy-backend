package main

import (
	"encoding/json"

	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

var HandleChat = func(config *requestConfig.APIGatewayWebsocketProxyRequest, message string, repo iface.ChatSessionRepository, client *router.Session) error {
	// marshall incoming payload
	var requestMessage request.Chat
	err := json.Unmarshal([]byte(message), &requestMessage)

	if err != nil {
		return err
	}

	sess, err := repo.GetEntityWithUsers(requestMessage.Session)

	if err != nil {
		return err
	}

	messageResponse, err := sess.AppendRequestMessage(requestMessage)

	if err != nil {
		return err
	}

	config.Logger = config.Logger.With(zap.String("id", messageResponse.ID), zap.String("sessionID", messageResponse.SessionID))

	body, err := json.Marshal(messageResponse)
	if err != nil {
		return err
	}

	logger := logging.LoggerFieldsFromEvent(config.Event)
	if logger.Error != nil {
		return logger.Error
	}

	err = client.Router.Send(&router.RouterSendInput{
		LoggerFields:             *logger.Fields,
		Source:                   "chat",
		Action:                   "chat",
		Procedure:                "send",
		Body:                     string(body),
		Recipients:               []string{sess.UserID()},
		OptOutGuaranteedDelivery: true,
	})

	if err != nil {
		return err
	}

	return nil
}
