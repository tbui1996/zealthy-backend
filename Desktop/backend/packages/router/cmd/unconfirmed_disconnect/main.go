//go:build !test
// +build !test

package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"context"
	"log"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Disconnecting client %s...", event.RequestContext.ConnectionID)
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, "unable to create session")
	}

	connectionId := event.RequestContext.ConnectionID
	email := event.RequestContext.Authorizer.(map[string]interface{})["email"].(string)

	svc := dynamodb.New(config.Session)

	err = Handler(UnconfirmedDisconnectRequest{
		ConnectionId: connectionId,
		Email:        email,
		Dynamo:       svc,
	})

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, err.Error())
	}

	log.Printf("Disconnected client %s.", event.RequestContext.ConnectionID)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
