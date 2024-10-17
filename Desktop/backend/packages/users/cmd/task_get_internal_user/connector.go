package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dto"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/repo"
	"go.uber.org/zap"
)

func connect(context context.Context, userID string) (dto.InternalUser, error) {
	result := &dto.InternalUser{}
	logger, err := logging.NewBasicLogger()

	if err != nil {
		return *result, err
	}

	logger = logger.With(zap.String("userID", userID))

	repo := repo.NewInternalUserRepository()

	logger.Debug("fetching user")
	result, err = repo.Find(userID)
	logger.Debug("user fetched")

	if err != nil {
		logger.Error(err.Error())
	}

	return *result, err
}

func main() {
	lambda.Start(connect)
}
