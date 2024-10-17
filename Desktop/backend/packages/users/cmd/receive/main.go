package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"go.uber.org/zap"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) (events.APIGatewayProxyResponse, error) {

	for _, message := range sqsEvent.Records {
		var loggerfields logging.LoggerFields
		logger, err := loggerfields.FromSQSMessage(message)
		if err != nil {
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{}, err
		}
		connectionID, ok := message.MessageAttributes["ConnectionId"]
		if !ok {
			logger.Error("couldn't get connection id")
		} else {
			logger.With(zap.String("connectionID", *connectionID.StringValue), zap.String("messageID", message.MessageId))
		}
	}
	return response.OK()
}

func main() {
	lambda.Start(HandleRequest)
}
