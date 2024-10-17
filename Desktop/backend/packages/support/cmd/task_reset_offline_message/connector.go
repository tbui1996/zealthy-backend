//go:build !test
// +build !test

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/repo"
	"go.uber.org/zap"
)

func connect(context context.Context, userID string) error {
	logger, err := logging.NewBasicLogger()

	if err != nil {
		return err
	}

	defer logging.SyncLogger(logger)

	logger = logger.With(zap.String("userId", userID))
	repo := repo.NewOfflineMessageNotificationRepoWithLogger(logger)

	return Handler(userID, HandlerDependencies{logger, repo})
}

func main() {
	lambda.Start(connect)
}
