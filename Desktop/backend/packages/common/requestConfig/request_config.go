package requestConfig

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"go.uber.org/zap"
)

type RequestConfig struct {
	Context context.Context
	Event   interface{}
	Logger  *zap.Logger
	Session *session.Session
}

func Must(config *RequestConfig, err *exception.SonarError) *RequestConfig {
	if err != nil {
		log.Fatalln(err.Error())
	}
	return config
}

// provisions a request wrapper that returns a config to be used for context, event, logger with parsed contextual information, and session
func NewRequestConfig(ctx context.Context, event interface{}) (*RequestConfig, *exception.SonarError) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if sess == nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "unable to get valid session. Please try again")
	}

	logger := logging.Must(logging.NewLoggerFromEvent(event))

	return &RequestConfig{
		Context: ctx,
		Event:   event,
		Logger:  logger,
		Session: sess,
	}, nil
}

// request wrapper with no logger. Used where a logger needs to be instatiated outside of the session
func NewRequestConfigNoLogger(ctx context.Context, event interface{}) (*RequestConfig, *exception.SonarError) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if sess == nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "unable to get valid session. Please try again")
	}

	return &RequestConfig{
		Context: ctx,
		Event:   event,
		Session: sess,
	}, nil
}

// provisions be used for context, parsing event data, logger, and session
func NewRequestConfigWithSessionConfig(ctx context.Context, event interface{}, config aws.Config) (*RequestConfig, *exception.SonarError) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            config,
	}))

	if sess == nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "unable to get valid session. Please try again")
	}

	logger := logging.Must(logging.NewLoggerFromEvent(event))

	return &RequestConfig{
		Context: ctx,
		Event:   event,
		Logger:  logger,
		Session: sess,
	}, nil
}
