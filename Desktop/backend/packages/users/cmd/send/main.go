package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/response"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	return response.OK()
}

func main() {
	lambda.Start(HandleRequest)
}
