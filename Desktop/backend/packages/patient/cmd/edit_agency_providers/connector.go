//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	var editAgencyProviderRequest request.EditAgencyProviderRequest
	if err := json.Unmarshal([]byte(config.Event.Body), &editAgencyProviderRequest); err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "unable to parse json body, cannot update patient"+err.Error())
	}

	agencyProviderId, ok := event.PathParameters["agency_provider_id"]
	if event.PathParameters == nil || !ok {
		errMsg := "parameter path {agency_provider_id} was not found"
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}

	editAgencyProviderRequest.AgencyProviderId = agencyProviderId

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "agency_provider_id was not valid")
	}

	repo, err := data.NewAgencyProviderRepository(config.Logger)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	deps := HandlerDeps{
		repo:   repo,
		logger: config.Logger,
	}

	return Handler(editAgencyProviderRequest, deps)
}

func main() {
	lambda.Start(connect)
}
