package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

func Handler(logger *zap.Logger, repo iface.ChatSessionRepository, userID string) ([]byte, *exception.SonarError) {
	sessions, err := repo.GetEntitiesByExternalID(userID)

	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "failed to get chat sessions for user. "+err.Error())
	}

	logger.Debug(fmt.Sprintf("received %d sessions, mapping to response dtos", len(sessions)))
	chatSessionResponseDTOs := make([]response.ChatSessionResponseDTO, len(sessions))
	for i := range sessions {
		chatSessionResponseDTOs[i] = sessions[i].ToResponseDTO()
	}

	body, err := json.Marshal(chatSessionResponseDTOs)
	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "failed to marshal sessions into JSON. "+err.Error())
	}

	return body, nil
}
