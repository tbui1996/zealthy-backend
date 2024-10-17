//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	filename, ok := event.QueryStringParameters["filename"]

	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "expected filename in query parameters")
	}

	bucket := os.Getenv("BUCKETNAME")
	s3Client := s3.New(config.Session)

	config.Logger.Info("pre_signed_upload_url started")
	key := uuid.New().String()
	res, err := Handler(PreSignedUploadUrlRequest{
		Logger:     config.Logger,
		S3API:      s3Client,
		BucketName: bucket,
		UniqueKey:  key,
		Filename:   filename,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}
	config.Logger.Info("pre_signed_upload_url complete")

	body, err := json.Marshal(response.PreSignedUrlResponse{
		URL: *res,
		Key: key,
	})

	if err != nil {
		errMsg := fmt.Sprintf("unable to marshal response. err (%s)", err)
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
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
	lambda.Start(HandleRequest)
}
