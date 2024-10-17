//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
)

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("pending_chat_session_create called")

	var createPendingChatSessionRequest model.PendingChatSessionCreate
	if err := json.Unmarshal([]byte(config.Event.Body), &createPendingChatSessionRequest); err != nil {
		errMsg := fmt.Errorf("unable to parse json body, cannot create pending chat session: %s", err.Error())
		config.Logger.Error(errMsg.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg.Error())
	}

	repo, err := chatSession.NewRDSChatSessionRepositoryWithSession(config.Session)
	db := dynamodb.New(config.Session)
	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String(os.Getenv("API_REGION")),
		Endpoint: aws.String(os.Getenv("WEBSOCKET_URL")),
	})

	if err != nil {
		errMessage := fmt.Sprintf("unable to open connection to db (%s)", err)
		config.Logger.Error(errMessage)
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, errMessage)
	}

	bodyBytes, err := handler(PendingSessionRequest{
		Request:  createPendingChatSessionRequest,
		Repo:     repo,
		Logger:   config.Logger,
		DynamoDB: db,
		API:      api,
	})

	if err != nil {
		config.Logger.Error("handling pending_chat_session_create: " + err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("pending_chat_session_create completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
