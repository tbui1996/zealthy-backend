//go:build !test
// +build !test

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

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

	config.Logger.Info("PUT /forms/{id}/close: started")

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

	date := time.Now()
	err = Handler(db, date, id)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not update closed date in database "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(connect)
}
