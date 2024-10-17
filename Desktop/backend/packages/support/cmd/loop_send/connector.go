//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
	"go.uber.org/zap"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("loop_send called")
	// marshall incoming payload
	var requestMessage request.Chat
	err = json.Unmarshal([]byte(event.Body), &requestMessage)
	if err != nil {
		config.Logger.Error("validating message" + err.Error())
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	config.Logger = config.Logger.With(zap.String("sessionID", requestMessage.Session))
	config.Logger.Debug("getting session")
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

	response, err := Handler(LoopSendRequest{
		Logger:   config.Logger,
		Repo:     repo,
		Message:  requestMessage,
		DynamoDB: db,
		API:      api,
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("loop_send completed, sending response")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(connect)
}
