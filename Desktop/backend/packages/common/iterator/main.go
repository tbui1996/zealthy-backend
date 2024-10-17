package iterator

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type IteratorHandler func(message events.SQSMessage) error

type Iterator struct {
	SQS sqsiface.SQSAPI

	QueueUrl *string

	handler IteratorHandler
}

func NewIteratorWithSession(sess *session.Session, queueName string) *Iterator {
	svc := sqs.New(sess)

	queue, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return nil
	}

	return &Iterator{
		SQS:      svc,
		QueueUrl: queue.QueueUrl,
	}
}

func (i *Iterator) UseHandler(f IteratorHandler) {
	i.handler = f
}

func (i *Iterator) Handle(event events.SQSEvent) error {
	// buffer channel so that none of the goroutines block
	eventErrors := make(chan error, len(event.Records))

	// run synchronously
	var wg sync.WaitGroup
	wg.Add(len(event.Records))
	for _, message := range event.Records {
		go func(closuredMessage events.SQSMessage) {
			defer wg.Done()

			err := i.handler(closuredMessage)

			if err != nil {
				eventErrors <- err
				return
			}

			// delete message if successful
			// ignoring error, there aren't any resolution steps
			// the code should be idempotent so that it can be rerun
			// if no errors occur then HandleRequest will return nil
			// and SQS will attempt to delete this message in the batch again
			// nolint errcheck
			i.SQS.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      i.QueueUrl,
				ReceiptHandle: aws.String(closuredMessage.ReceiptHandle),
			})
		}(message)
	}

	// wait for all goroutines to complete
	wg.Wait()

	// indicate that no more data will come through channel
	close(eventErrors)

	// we don't know if there's any errors preemptively
	// nolint prealloc
	var errors []string
	for err := range eventErrors {
		errors = append(errors, err.Error())
	}

	// return an error so that SQS does not delete remaining messages from the queue
	if len(errors) > 0 {
		errMessage := strings.Join(errors, ", ")
		return fmt.Errorf(errMessage)
	}

	// SQS will delete any remaining messages in this batch from the queue
	// This should be none, but in case SQS.DeleteMessage failed, it's a second
	// delete attempt, this time performed by AWS
	return nil
}
