package requestConfig

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type APIGatewayWebsocketProxyRequest struct {
	Context context.Context
	Event   events.APIGatewayWebsocketProxyRequest
	Logger  *zap.Logger
	Session *session.Session
}

func (rc *RequestConfig) ToAPIGatewayWebsocketProxyRequest() (*APIGatewayWebsocketProxyRequest, error) {
	event, ok := rc.Event.(events.APIGatewayWebsocketProxyRequest)
	if !ok {
		return nil, errors.New("type of event is not events.APIGatewayWebsocketProxyRequest")
	}

	return &APIGatewayWebsocketProxyRequest{
		Context: rc.Context,
		Event:   event,
		Logger:  rc.Logger,
		Session: rc.Session,
	}, nil
}
