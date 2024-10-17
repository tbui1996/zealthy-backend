package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/iface"
	"go.uber.org/zap"
)

type HandlerDependencies struct {
	logger *zap.Logger
	repo   iface.OfflineMessageNotificationRepo
}

func Handler(userID string, deps HandlerDependencies) error {
	deps.logger.Info("resetting message")

	err := deps.repo.Remove(userID)

	if err != nil {
		deps.logger.Error(fmt.Sprintf("failed to remove notification: %s", err.Error()))
	}

	return err
}
