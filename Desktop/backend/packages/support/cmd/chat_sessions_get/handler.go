package main

import (
	"encoding/json"
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

func Handler(userID string, group string, logger *zap.Logger, repo iface.ChatSessionRepository) ([]byte, error) {
	chatType, ok := model.SupportTypes[group]
	if !ok {
		errMsg := fmt.Errorf("expected a valid chat type for group %s", group)
		logger.Error(errMsg.Error())
		return nil, errMsg
	}

	chatSessions, err := repo.GetEntities(userID, model.ChatTypeToString(chatType))
	if err != nil {
		errMsg := fmt.Errorf("unable to get connection items %s", err.Error())
		logger.Error(errMsg.Error())
		return nil, errMsg
	}

	logger.Debug(fmt.Sprintf("received %d sessions, mapping to response dtos", len(chatSessions)))
	chatSessionResponseDTOs := make([]response.ChatSessionResponseDTO, len(chatSessions))
	for i := range chatSessions {
		chatSessionResponseDTOs[i] = chatSessions[i].ToResponseDTO()
	}

	body, err := json.Marshal(chatSessionResponseDTOs)
	if err != nil {
		errMsg := fmt.Errorf("unable to marshal connected chatSessions response %s", err.Error())
		logger.Error(errMsg.Error())
		return nil, errMsg
	}

	return body, nil
}
