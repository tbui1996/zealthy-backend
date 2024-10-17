//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}
	defer logging.SyncLogger(config.Logger)

	repo, err := data.NewAppointmentRepository(config.Logger)

	if err != nil {
		config.Logger.Error(err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "Unable to connect to database")
	}

	appointments, err := Handler(HandlerDeps{
		repo:   repo,
		logger: config.Logger,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	body, err := json.Marshal(response.ResultWrapper{
		Result: *appointments,
	})

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while marshalling: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "There was an error marshaling the response from the database.")
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
