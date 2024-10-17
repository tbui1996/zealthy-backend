package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

func Handler(logger *zap.Logger, repo iface.ChatSessionRepository, request request.AssignPendingChatSessionRequestInternal) ([]byte, *exception.SonarError) {
	logger.Debug(fmt.Sprintf("assigning session to %s", request.InternalUserID))

	sessionID, _ := strconv.Atoi(request.SessionID)
	logger = logger.With(zap.Int("sessionID", sessionID))

	sess, err := repo.AssignPending(sessionID, request.InternalUserID)
	if err != nil {
		errMessage := fmt.Sprintf("error storing item: (%s)", err)
		logger.Error(errMessage)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	// get messages to hydrate
	logger.Debug("getting messages for session")
	messages, err := sess.GetMessages()
	if err != nil {
		errMessage := fmt.Sprintf("error getting pending chat messages: %+v (%s)", sessionID, err)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	logger.Debug("received messages for session, marshalling")
	body, err := json.Marshal(messages)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "unable to marshal pending chat messages response "+err.Error())
	}

	return body, nil
}
