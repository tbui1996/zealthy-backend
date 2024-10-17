package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

func Handler(logger *zap.Logger, repo iface.ChatMessageRepository, sessionID string) ([]byte, *exception.SonarError) {
	messages, err := repo.GetMessagesForSession(sessionID)
	if err != nil {
		errMessage := fmt.Sprintf("error getting chat messages from dynamo: %+v (%s)", sessionID, err)
		logger.Error(errMessage)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	logger.Debug("received messages for session, marshalling body")
	body, err := json.Marshal(messages)
	if err != nil {
		logger.Error("unable to marshal chat messages response " + err.Error())
		return nil, exception.NewSonarError(http.StatusBadRequest, "unable to marshal chat messages response "+err.Error())
	}

	return body, nil
}
