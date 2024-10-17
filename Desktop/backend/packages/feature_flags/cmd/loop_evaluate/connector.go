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
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	featureflags "github.com/circulohealth/sonar-backend/packages/common/feature_flags"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/response"
)

func Connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	db, err := dao.OpenConnectionWithTablePrefix(dao.FeatureFlags)

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while creating db: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	evaluator := featureflags.NewEvaluatorWithDB(db)

	results, err := evaluator.Evaluate()

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while evaluating: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	body, err := json.Marshal(response.ResultWrapper{
		Result: results.Map(),
	})

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while marshalling: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "Error while marshalling")
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
	lambda.Start(Connect)
}
