//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("GET /cloud/get_file: started")

	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config")
	}

	defer logging.SyncLogger(config.Logger)

	db, err := dao.OpenConnectionWithTablePrefix(dao.Cloud)
	if err != nil {
		errMSg := "unable to open connection to db"
		config.Logger.Error(errMSg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMSg)
	}

	files, err := handler(db)
	if err != nil {
		errMSg := "unable to get files from DB"
		config.Logger.Error(errMSg + fmt.Sprintf(" %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMSg)
	}

	body, err := json.Marshal(files)
	if err != nil {
		errMsg := "unable to marshal the response from the database"
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMsg)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
