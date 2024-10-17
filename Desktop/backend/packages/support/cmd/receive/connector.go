//go:build !test
// +build !test

package main

import (
	"context"
	"fmt"

	chatSession "github.com/circulohealth/sonar-backend/packages/support/pkg/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/circulohealth/sonar-backend/packages/common/iterator"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	repo, err := chatSession.NewRDSChatSessionRepositoryWithSession(sess)
	if err != nil {
		return fmt.Errorf("unable to open connection to DB (%s)", err.Error())
	}

	sqsIterator := iterator.NewIteratorWithSession(sess, "sonar-service-support-receive")

	sqsIterator.UseHandler(func(message events.SQSMessage) error {
		return Handler(message, repo)
	})

	return sqsIterator.Handle(event)
}

func main() {
	lambda.Start(HandleRequest)
}
