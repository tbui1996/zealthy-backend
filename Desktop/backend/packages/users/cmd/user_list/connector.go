//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/response"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("list_users called")

	userPoolID := os.Getenv("USER_POOL_ID")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not connect to database, unable to retrieve form "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	idp := cognitoidentityprovider.New(config.Session)

	registry := mapper.NewRegistry(&mapper.NewRegistryInput{
		DB:         db,
		IDP:        idp,
		UserPoolId: &userPoolID,
		Logger:     config.Logger,
	})

	output, err := handler(HandlerInput{
		Logger:                config.Logger,
		QueryStringParameters: event.QueryStringParameters,
		Registry:              registry,
	})

	if err != nil {
		config.Logger.Error("users_list failed: " + err.Error())
		return exception.NewSonarError(http.StatusInternalServerError, err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	wrapped := response.Users{
		Users: output,
	}

	body, err := json.Marshal(wrapped)

	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "unable to marshal users response "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("list_users complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
