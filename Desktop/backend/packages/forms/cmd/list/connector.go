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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	config.Logger.Info("GET /forms: started")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to open connection to db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	forms, formsErr := getAllForms(db)

	if formsErr != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to get forms from database "+formsErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	body, err := json.Marshal(forms)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to marshal forms response "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("GET /forms: complete")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
