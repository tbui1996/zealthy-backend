package main

import (
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

func Handler(input request.DeleteFlagRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	if input.FlagId == 0 {
		errMsg := "ID is required"
		deps.logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}
	flag, err := deps.repo.Find(input.FlagId)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	newError := deps.repo.Delete(flag)

	if newError != nil {
		return newError.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
