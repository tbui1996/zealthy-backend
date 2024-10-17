package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/response"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.FeatureFlagRepository
	logger *zap.Logger
}

func Handler(deps HandlerDeps) (*[]response.FeatureFlagResponse, error) {
	flags, err := deps.repo.FindAll()

	if err != nil {
		deps.logger.Error(err.Error())
		return nil, fmt.Errorf("unable to find records in database")
	}

	responses := make([]response.FeatureFlagResponse, len(*flags))

	for i := range *flags {
		flag := (*flags)[i]
		responses[i] = response.FeatureFlagResponse{
			Id:        flag.Id,
			Key:       flag.Key,
			Name:      flag.Name,
			CreatedAt: flag.CreatedAt.String(),
			CreatedBy: *flag.CreatedBy,
			UpdatedAt: flag.UpdatedAt.String(),
			UpdatedBy: *flag.UpdatedBy,
			IsEnabled: flag.IsEnabled,
		}
	}

	return &responses, nil
}
