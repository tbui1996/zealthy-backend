package requestConfig

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type APIGatewayProxyRequest struct {
	Context context.Context
	Event   events.APIGatewayProxyRequest
	Logger  *zap.Logger
	Session *session.Session
}

func (rc *RequestConfig) ToAPIGatewayProxyRequest() (*APIGatewayProxyRequest, error) {
	event, ok := rc.Event.(events.APIGatewayProxyRequest)
	if !ok {
		return nil, errors.New("type of event is not events.APIGatewayProxyRequest")
	}

	return &APIGatewayProxyRequest{
		Context: rc.Context,
		Event:   event,
		Logger:  rc.Logger,
		Session: rc.Session,
	}, nil
}
