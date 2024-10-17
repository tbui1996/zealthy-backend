package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.FeatureFlagRepository
	logger *zap.Logger
}

func Handler(input request.CreateFlagRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {

	if input.Key == "" || input.Name == "" {
		deps.logger.Error("FlagKey and Name are both required")
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "FlagKey and Name are both required")
	}

	flag := model.NewFeatureFlagWithUserId(model.FeatureFlag{
		Key:  input.Key,
		Name: input.Name,
	}, &input.UserId)

	err := deps.repo.Save(flag)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
