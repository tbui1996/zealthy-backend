package exception

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func ErrorMessage(statusCode int, errorMessage string) (resp events.APIGatewayProxyResponse, err error) {
	log.Println(errorMessage)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       errorMessage,
	}, nil
}

func ErrorMessageApiGatewayV2(statusCode int, errorMessage string) (events.APIGatewayV2HTTPResponse, error) {
	log.Println(errorMessage)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body:       errorMessage,
	}, nil
}
