package router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
)

type Config struct {
	SendQueueName    string
	ReceiveQueueName string
}

type RouterSendInput struct {
	LoggerFields logging.LoggerFields
	Action       string
	Procedure    string
	Recipients   []string
	Source       string
	Body         string

	/* Optional */
	OptOutGuaranteedDelivery bool
}

type Client struct {
	client *sqs.SQS
	config *Config
}

// New TODO: Replace all instances with code in session.go
func New(client *sqs.SQS, config *Config) *Client {
	rc := &Client{
		client: client,
		config: config,
	}

	return rc
}

func validateInput(input *RouterSendInput) error {
	if input.Action == "" {
		return fmt.Errorf("expected Action to be set in input")
	}

	if input.Source == "" {
		return fmt.Errorf("expected Source to be set in input")
	}

	if input.Procedure == "" {
		return fmt.Errorf("expected Procedure to be set in input")
	}

	if input.Recipients == nil {
		return fmt.Errorf("expected Recipients to be set in input")
	}

	if err := input.LoggerFields.Validate(); err != nil {
		return err
	}

	return nil
}

func (rc *Client) Send(input *RouterSendInput) error {
	err := validateInput(input)

	if err != nil {
		return err
	}

	queue, err := rc.client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(rc.config.SendQueueName),
	})

	if err != nil {
		return err
	}

	loggerFieldsBytes, err := json.Marshal(input.LoggerFields)
	if err != nil {
		return err
	}

	messageAttributes := map[string]*sqs.MessageAttributeValue{
		"Source": {
			DataType:    aws.String("String"),
			StringValue: aws.String(input.Source),
		},
		"Procedure": {
			DataType:    aws.String("String"),
			StringValue: aws.String(input.Procedure),
		},
		"Action": {
			DataType:    aws.String("String"),
			StringValue: aws.String(input.Action),
		},
		"LoggerFields": {
			DataType:    aws.String("String"),
			StringValue: aws.String(string(loggerFieldsBytes)),
		},
	}

	if len(input.Recipients) > 0 {
		messageAttributes["Recipients"] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(strings.Join(input.Recipients, ",")),
		}
	}

	if input.OptOutGuaranteedDelivery {
		messageAttributes["OptOutGuaranteedDelivery"] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("true"),
		}
	}

	// Send message through SQS to dispatch
	_, err = rc.client.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: messageAttributes,
		MessageBody:       aws.String(input.Body),
		QueueUrl:          queue.QueueUrl,
	})

	return err
}
