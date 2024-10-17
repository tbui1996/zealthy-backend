//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()

	if err != nil {
		errMsg := fmt.Errorf("error setting up config: %s", err.Error())
		log.Print(errMsg.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMsg.Error())
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("update_user called")

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)

	if err != nil {
		return errorResponse(err, "could not connect to database", config.Logger)
	}

	idp := cognitoidentityprovider.New(config.Session)

	userPoolID := os.Getenv("USER_POOL_ID")
	registry := mapper.NewRegistry(&mapper.NewRegistryInput{DB: db, Logger: config.Logger, IDP: idp, UserPoolId: &userPoolID})

	var req request.UpdateUserRequest
	err = json.Unmarshal([]byte(config.Event.Body), &req)
	if err != nil {
		return errorResponse(err, "invalid request body", config.Logger)
	}

	webSocketUrl := os.Getenv("WEBSOCKET_URL")
	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String(webSocketUrl),
	})
	usersDB := dynamodb.New(config.Session)

	err = handler(req, UpdateUserDependencies{Registry: registry, Logger: config.Logger, Api: api, Db: usersDB})

	if err != nil {
		return errorResponse(err, "unable to update user", config.Logger)
	}

	config.Logger.Info("update_user complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
