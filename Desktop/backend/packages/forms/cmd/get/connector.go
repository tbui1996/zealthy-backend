//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/repo"
)

// HandleRequest event form.GetForm -> *form.Form
func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("Get /forms/{id}: started")

	if event.PathParameters == nil {
		return exception.NewSonarError(http.StatusBadRequest, "expected path parameters to exist").ToAPIGatewayV2HTTPResponse(), nil
	}

	id, ok := event.PathParameters["id"]
	if !ok {
		return exception.NewSonarError(http.StatusBadRequest, "parameter path {id} was not found").ToAPIGatewayV2HTTPResponse(), nil
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not connect to database, unable to retrieve form "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	repository := repo.NewRepository(db)
	inputs, err := repository.Inputs(id)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not retrieve form inputs from db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	form, err := repository.Form(id)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not retrieve form from db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	body, err := json.Marshal(response.GetFormResponse{Form: *form, Inputs: inputs})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to parse JSON response "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("Get /forms/{id}: complete")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
