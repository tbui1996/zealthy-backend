//go:build !test
// +build !test

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"go.uber.org/zap"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	fileID, ok := event.PathParameters["id"]
	config.Logger = config.Logger.With(zap.String("fileID", fileID))
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "File ID not provided in request. Please try again.")
	}

	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("file_download called")

	bucket := os.Getenv("BUCKETNAME")
	svc := s3.New(config.Session)

	resp, err := handler(FileDownloadRequest{
		S3:         svc,
		BucketName: bucket,
		FileId:     fileID,
		Logger:     config.Logger,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("file_download completed")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Location": *resp,
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
