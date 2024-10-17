package requestConfig

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type APIGatewayV2HTTPRequest struct {
	Context context.Context
	Event   events.APIGatewayV2HTTPRequest
	Logger  *zap.Logger
	Session *session.Session
}

func (rc *RequestConfig) ToAPIGatewayV2HTTPRequest() (*APIGatewayV2HTTPRequest, error) {
	event, ok := rc.Event.(events.APIGatewayV2HTTPRequest)
	if !ok {
		return nil, errors.New("type of event is not events.APIGatewayV2HTTPRequest")
	}

	return &APIGatewayV2HTTPRequest{
		Context: rc.Context,
		Event:   event,
		Logger:  rc.Logger,
		Session: rc.Session,
	}, nil
}
