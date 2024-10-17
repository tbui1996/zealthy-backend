//go:build !test
// +build !test

package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

	"context"
	"log"
	"net/http"

	commonEvents "github.com/circulohealth/sonar-backend/packages/common/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayWebsocketProxyRequest()

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, "unable to get configuration for connect request")
	}

	userID := event.RequestContext.Authorizer.(map[string]interface{})["userID"].(string)
	connectionID := event.RequestContext.ConnectionID
	log.Printf("Connecting client: %s", connectionID)
	if len(connectionID) == 0 {
		return exception.ErrorMessage(http.StatusBadRequest, "Error marshaling connection item")
	}

	db := dynamodb.New(config.Session)
	eventBridge := eventbridge.New(config.Session)

	err = Handler(ConnectRequest{
		UserID:       userID,
		ConnectionId: connectionID,
		Dynamo:       db,
		EventPublisher: &commonEvents.EventPublisher{
			EventBridge: eventBridge,
		},
	})

	if err != nil {
		return exception.ErrorMessage(http.StatusBadRequest, fmt.Sprintf("Failed to connect client %s. Error on PutItem.", connectionID))
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
