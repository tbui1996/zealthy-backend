package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.FeatureFlagRepository
	logger *zap.Logger
}

func Handler(input request.PatchFlagRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	if (input.Name == nil || *input.Name == "") && input.IsEnabled == nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Either Name or IsEnabled are required")
	}

	flag, err := deps.repo.Find(input.FlagId)

	if err != nil {
		deps.logger.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
		return err.ToSonarApiGatewayResponse()
	}

	if input.Name != nil && *input.Name != "" {
		flag.Name = *input.Name
	}

	if input.IsEnabled != nil {
		flag.IsEnabled = *input.IsEnabled
	}

	err = deps.repo.Save(flag)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
