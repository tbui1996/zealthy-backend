//go:build !test
// +build !test

package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	userId, ok := config.Event.RequestContext.Authorizer.Lambda["userID"].(string)
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "user id not found")
	}

	var deleteFlagRequest request.DeleteFlagRequest

	flagId, ok := event.PathParameters["id"]
	if event.PathParameters == nil || !ok {
		errMsg := "parameter path {id} was not found"
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}

	deleteFlagRequest.FlagId, err = strconv.Atoi(flagId)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "id was not a number")
	}

	repo, err := data.NewFeatureFlagRepository(userId, config.Logger)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	deps := HandlerDeps{
		repo:   repo,
		logger: config.Logger,
	}

	return Handler(deleteFlagRequest, deps)
}

func main() {
	lambda.Start(connect)
}
