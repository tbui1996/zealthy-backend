//go:build !test
// +build !test

package main

import (
	"context"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/patients"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"go.uber.org/zap"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	logging.SyncLogger(config.Logger)

	config.Logger.Info("patients_get called")

	userId, ok := config.Event.RequestContext.Authorizer.Lambda["userID"].(string)
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "user id not found")
	}

	config.Logger = config.Logger.With(zap.String("userId", userId))
	config.Logger.Debug(fmt.Sprintf("getting patients for user with id %s", userId))

	repo, err := patients.NewPatientRepository()
	if err != nil {
		config.Logger.Error(err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	body, err := Handler(PatientsGetRequest{Logger: config.Logger, Repo: repo, UserId: userId})
	if err != nil {
		config.Logger.Error(err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("patients_get completed")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(connector)
}
