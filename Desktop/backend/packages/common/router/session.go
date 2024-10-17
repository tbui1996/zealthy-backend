package router

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Session struct {
	Router Router
}

func NewClientWithConfig(config *Config) *Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	r := New(svc, config)

	return &Session{
		Router: r,
	}
}

func NewClientWithConfigWithSession(config *Config, session *session.Session) *Session {
	svc := sqs.New(session)

	r := New(svc, config)

	return &Session{
		Router: r,
	}
}
