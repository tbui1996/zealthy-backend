//go:build !test
// +build !test

package main

import (
	"context"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()
	if err != nil {
		return exception.ErrorMessage(http.StatusInternalServerError, "error setting up config")
	}
	logging.SyncLogger(config.Logger)

	config.Logger.Info("support send called")

	repo, err := session.NewRDSChatSessionRepositoryWithSession(config.Session)

	if err != nil {
		return exception.ErrorMessage(http.StatusInternalServerError, fmt.Sprintf("unable to open connection to DB (%s)", err.Error()))
	}

	client := router.NewClientWithConfigWithSession(&router.Config{
		SendQueueName:    "sonar-service-support-send",
		ReceiveQueueName: "sonar-service-support-receive",
	}, config.Session)

	response, sErr := Handler(config, repo, client)
	if sErr != nil {
		config.Logger.Error("support send fails: " + sErr.Error())
		return sErr.ToAPIGatewayProxyResponse(), nil
	}

	config.Logger.Info("support send completed, sending response")

	return *response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
