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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/response"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	db := dynamodb.New(config.Session)

	deps := HandlerDeps{
		db,
		config.Logger,
	}

	results, err := Handler(deps)

	if err != nil {
		config.Logger.Error(fmt.Sprintf("error in handler: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "internal error")
	}

	resultBytes, err := json.Marshal(response.ResultWrapper{
		Result: results,
	})

	if err != nil {
		config.Logger.Error(fmt.Sprintf("error marshalling results: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error marshalling online users results")
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(resultBytes),
	}, nil
}

func main() {
	lambda.Start(connect)
}
