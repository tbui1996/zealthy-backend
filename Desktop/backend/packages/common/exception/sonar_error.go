package exception

import (
	"github.com/aws/aws-lambda-go/events"
)

type SonarError struct {
	StatusCode int
	message    string
}

func NewSonarError(statusCode int, errorMessage string) *SonarError {
	return &SonarError{
		message:    errorMessage,
		StatusCode: statusCode,
	}
}

func (e *SonarError) Error() string {
	return e.message
}

func (e *SonarError) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: e.StatusCode,
		Body:       e.Error(),
	}
}

func (e *SonarError) ToAPIGatewayV2HTTPResponse() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: e.StatusCode,
		Body:       e.Error(),
	}
}
