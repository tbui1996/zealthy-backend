//go:build !test
// +build !test

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"
	"go.uber.org/zap"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	logging.SyncLogger(config.Logger)

	config.Logger.Info("chat_messages_get called")

	if event.PathParameters == nil {
		config.Logger.Error("expected path parameters to exist")
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "expected path parameters to exist")
	}

	sessionID, ok := event.PathParameters["id"]
	if !ok {
		config.Logger.Error("parameter path {id} was not found")
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "parameter path {id} was not found")
	}

	config.Logger = config.Logger.With(zap.String("sessionID", sessionID))
	config.Logger.Debug("getting chat messages for session")

	db := dynamodb.New(config.Session)
	repo := chatSession.NewDynamoDBChatMessageRepositoryWithDB(db)

	body, sErr := Handler(config.Logger, repo, sessionID)

	if sErr != nil {
		config.Logger.Error(err.Error())
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("chat_messages_get completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(connector)
}
