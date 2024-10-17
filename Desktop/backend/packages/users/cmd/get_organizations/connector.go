//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()

	if err != nil {
		errMsg := fmt.Errorf("error setting up config: %s", err.Error())
		log.Print(errMsg.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMsg.Error())
	}

	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("get_organizations called")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)

	if err != nil {
		return errorResponse(err, "could not connect to database, unable to retrieve form", config.Logger)
	}

	registry := mapper.NewRegistry(&mapper.NewRegistryInput{DB: db, Logger: config.Logger, IDP: nil, UserPoolId: nil})

	output, err := handler(GetOrganizationsRequest{Registry: registry, Logger: config.Logger})

	if err != nil {
		return errorResponse(err, "unable to get organizations", config.Logger)
	}

	body, err := json.Marshal(output)

	if err != nil {
		return errorResponse(err, "unable to marshal organizations response", config.Logger)
	}

	config.Logger.Info("get_organizations complete, sending response")

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
