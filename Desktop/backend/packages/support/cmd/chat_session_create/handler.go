package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

func Handler(logger *zap.Logger, repo iface.ChatSessionRepository, createRequest request.ChatSessionRequestInternal) ([]byte, *exception.SonarError) {
	chatSessionCreateRequest := &request.ChatSessionCreateRequest{
		InternalUserID: &createRequest.InternalUserID,
		UserID:         createRequest.UserID,
		ChatOpen:       true,
		Topic:          "",
		Created:        time.Now().Unix(),
	}

	logger.Debug("storing chat session")
	sess, err := repo.Create(chatSessionCreateRequest)

	if err != nil {
		errMessage := fmt.Sprintf("error storing item: %+v (%s)", chatSessionCreateRequest, err)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	if sess == nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "expected session to exist after creating")
	}

	logger = logger.With(zap.String("sessionID", sess.ID()))
	logger = logger.With(zap.String("internalUserID", createRequest.InternalUserID))
	logger = logger.With(zap.String("externalUserID", createRequest.UserID))
	logger.Debug("stored chat session successfully")

	chatSessionResponse := response.ChatSessionResponse{
		ID: sess.ID(),
	}

	bodyBytes, err := json.Marshal(chatSessionResponse)
	if err != nil {
		errMessage := fmt.Sprintf("error marshalling session id response: %s", err)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return bodyBytes, nil
}
