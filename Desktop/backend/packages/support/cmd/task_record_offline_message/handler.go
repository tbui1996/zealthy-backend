package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/dto"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/iface"
	"go.uber.org/zap"
)

type HandlerDependencies struct {
	logger *zap.Logger
	repo   iface.OfflineMessageNotificationRepo
}

func Handler(userID string, deps HandlerDependencies) (dto.RecordOfflineMessageDTO, error) {
	created, err := deps.repo.Create(userID)

	if err != nil {
		deps.logger.Error(fmt.Sprintf("unable to create notification: %s", err.Error()))
	}

	return dto.RecordOfflineMessageDTO{
		Created: created,
	}, err
}
