//go:build !test
// +build !test

package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

	"context"
	"log"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, fmt.Sprintf("unable to get config for session (%s)", err))
	}

	email := event.RequestContext.Authorizer.(map[string]interface{})["email"].(string)
	connectionID := event.RequestContext.ConnectionID
	log.Printf("Connecting client: %s", connectionID)
	if len(connectionID) == 0 {
		return exception.ErrorMessage(http.StatusBadRequest, "Error marshaling connection item")
	}

	svc := dynamodb.New(config.Session)

	err = Handler(UnconfirmedConnectRequest{
		Email:        email,
		ConnectionId: connectionID,
		Dynamo:       svc,
	})

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, err.Error())
	}

	log.Printf("Connected client %s.", connectionID)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
