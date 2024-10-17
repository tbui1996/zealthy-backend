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
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"log"
	"os"

	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/forward"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, fmt.Sprintf("unable to get configuration for session (%s)", err.Error()))
	}

	var forwarderBroadcastDTO forward.ForwarderBroadcastDTO
	err = json.Unmarshal([]byte(event.Body), &forwarderBroadcastDTO)
	if err != nil {
		log.Printf("Error unmarshalling request data: %s\n%s", event.Body, err)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Error unmarshalling request data")
	}

	svc := dynamodb.New(config.Session)

	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String(os.Getenv("WEBSOCKET_URL")),
	})

	forwarderBroadcastDTO.ApiGatewayManagementApi = api
	forwarderBroadcastDTO.DynamoDB = svc
	forwarderBroadcastDTO.Logger = config.Logger

	err = Handler(forwarderBroadcastDTO)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	return response.OKv2()
}

func main() {
	lambda.Start(HandleRequest)
}
