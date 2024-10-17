//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/request"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"go.uber.org/zap"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config").ToAPIGatewayV2HTTPResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	userId, ok := config.Event.RequestContext.Authorizer.Lambda["userID"].(string)
	if !ok {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "user id not found")
	}

	body := []byte(config.Event.Body)
	if config.Event.IsBase64Encoded {
		if body, err = base64.StdEncoding.DecodeString(config.Event.Body); err != nil {
			return events.APIGatewayV2HTTPResponse{}, err
		}
	}

	var fileUploadRequest request.FileUploadRequest
	if err := json.Unmarshal(body, &fileUploadRequest); err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, fmt.Sprintf("unable to parse json body, cannot create form (%s)", err.Error()))
	}

	config.Logger = config.Logger.With(zap.String("userID", userId))
	config.Logger = config.Logger.With(zap.String("sessionID", fileUploadRequest.ChatId))
	config.Logger.Info("file_upload called")

	db, dbConnErr := dao.OpenConnectionWithTablePrefix(dao.Cloud)

	if dbConnErr != nil {
		return exception.NewSonarError(http.StatusServiceUnavailable, "unable to connect to database: "+dbConnErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	svc := s3.New(config.Session)

	handlerInput := &FileUploadHandler{
		DB:            db,
		BucketName:    os.Getenv("BUCKETNAME"),
		Logger:        config.Logger,
		S3:            svc,
		Username:      userId,
		UploadRequest: fileUploadRequest,
	}

	fileResponse, err := handler(handlerInput)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get valid file response (%s)", err.Error())
		config.Logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}
	config.Logger.Info("file_upload complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(fileResponse),
	}, nil
}

func main() {
	lambda.Start(connector)
}
