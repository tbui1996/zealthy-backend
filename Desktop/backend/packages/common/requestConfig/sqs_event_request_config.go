package requestConfig

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type SQSEventRequest struct {
	Context context.Context
	Event   events.SQSEvent
	Logger  *zap.Logger
	Session *session.Session
}

func (rc *RequestConfig) ToSQSEventRequest() (*SQSEventRequest, error) {
	event, ok := rc.Event.(events.SQSEvent)
	if !ok {
		return nil, errors.New("type of event is not events.APIGatewayWebsocketProxyRequest")
	}

	return &SQSEventRequest{
		Context: rc.Context,
		Event:   event,
		Logger:  rc.Logger,
		Session: rc.Session,
	}, nil
}
