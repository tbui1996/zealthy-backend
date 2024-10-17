package main

import (
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
)

type TypingAction struct {
	Action    string `json:"action"`
	SessionId string `json:"sessionId"`
}

var HandleTyping = func(config *requestConfig.APIGatewayWebsocketProxyRequest, message string, repo iface.ChatSessionRepository, client *router.Session) error {
	var req request.TypingActionRequest
	err := json.Unmarshal([]byte(message), &req)

	verb := "started"
	if req.Action == "stop" {
		verb = "stopped"
	}

	config.Logger.Debug(fmt.Sprintf("user %s has %sed typing for session %s", req.UserID, verb, req.SessionID))

	if err != nil {
		return err
	}

	sess, err := repo.GetEntityWithUsers(req.SessionID)

	if err != nil {
		return err
	}

	logger := logging.LoggerFieldsFromEvent(config.Event)
	if logger.Error != nil {
		return logger.Error
	}

	msg, err := json.Marshal(TypingAction{Action: req.Action, SessionId: req.SessionID})
	if err != nil {
		return err
	}

	config.Logger.Debug(fmt.Sprintf("Sending typing action message %s", string(msg)))

	err = client.Router.Send(&router.RouterSendInput{
		LoggerFields:             *logger.Fields,
		Source:                   "chat",
		Action:                   "chat",
		Procedure:                "typing",
		Body:                     string(msg),
		Recipients:               []string{sess.UserID()},
		OptOutGuaranteedDelivery: true,
	})

	if err != nil {
		return err
	}

	config.Logger.Debug("typing action message sent")

	return nil
}
