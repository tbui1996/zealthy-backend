//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("PUT /cloud/delete_file/{id}: started")

	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config")
	}

	defer logging.SyncLogger(config.Logger)

	var deleteFileRequest request.DeleteFile
	err = json.Unmarshal([]byte(event.Body), &deleteFileRequest)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Request body not accepted. Please try again."+err.Error())
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Cloud)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "unable to connect to db "+err.Error())
	}

	date := time.Now()
	result := handler(&AddDeleteDateInput{
		Date:   &date,
		Db:     db,
		FileID: deleteFileRequest.ID,
	})

	if result != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "unable to add deleted_at date to form item "+result.Error())
	}

	log.Println("file_delete complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "Successfully deleted file.",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
