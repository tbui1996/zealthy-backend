//go:build !test
// +build !test

package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config").ToAPIGatewayV2HTTPResponse(), nil
	}
	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("associate_file called")
	file, err := parseRequest(event)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Request body not accepted. Please try again.")
	}

	sonarDb, err := dao.OpenConnectionWithTablePrefix(dao.Cloud)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Unable to find valid session. Please try again.")
	}

	config.Logger.Info("connecting to doppler")
	dopplerDb, err := dao.OpenConnectionToDoppler()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	err = associateFile(&AssociateFileHandler{
		DopplerDb: dopplerDb,
		SonarDb:   sonarDb,
		File:      file,
		Logger:    config.Logger,
	})
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "Successfully updated memberId.",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
