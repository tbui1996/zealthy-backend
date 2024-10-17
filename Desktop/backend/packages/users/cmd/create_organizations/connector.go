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
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("POST /users/organizations started")
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()

	if err != nil {
		errMsg := fmt.Errorf("error setting up config: %s", err.Error())
		log.Print(errMsg.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMsg.Error())
	}
	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("POST /users/organizations called")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)

	if err != nil {
		return errorResponse(err, "could not connect to database, unable to retrieve form", config.Logger)
	}
	registry := mapper.NewRegistry(&mapper.NewRegistryInput{DB: db, Logger: config.Logger, IDP: nil, UserPoolId: nil})

	var req request.CreateOrganizationRequest
	err = json.Unmarshal([]byte(config.Event.Body), &req)
	if err != nil {
		return errorResponse(err, "invalid request body", config.Logger)
	}

	output, err := handler(CreateOrganizationsRequest{Registry: registry, Logger: config.Logger, OrganizationName: req.Name})

	if err != nil {
		return errorResponse(err, "unable to create organizations", config.Logger)
	}

	body, err := json.Marshal(output)

	if err != nil {
		return errorResponse(err, "unable to marshal organizations response", config.Logger)
	}

	config.Logger.Info("create_organizations complete, sending response")

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
