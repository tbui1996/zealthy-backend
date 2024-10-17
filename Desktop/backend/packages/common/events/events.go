package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
)

type SenderType string

const (
	InternalSenderType SenderType = "internal"
	LoopSenderType     SenderType = "loop"
)

type MessageSentEvent struct {
	ReceiverId string
	SenderId   string
	SenderType SenderType
	SentAt     int64
	SessionId  string
	MessageId  string
}

type ConnectionCreatedEvent struct {
	UserID    string
	CreatedAt int64
}

type EventPublisher struct {
	EventBridge eventbridgeiface.EventBridgeAPI
}

func (pub *EventPublisher) PublishConnectionCreatedEvent(userId string, service string) error {
	event := ConnectionCreatedEvent{
		userId,
		time.Now().Unix(),
	}

	encodedEvent, err := json.Marshal(event)

	if err != nil {
		return err
	}

	eventBridgeInput := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventconstants.SERVICE_EVENT_BUS),
				Detail:       aws.String(string(encodedEvent)),
				DetailType:   aws.String(eventconstants.CONNECTION_CREATED_EVENT),
				Resources:    []*string{},
				Source:       aws.String(service),
			},
		},
	}

	res, err := pub.EventBridge.PutEvents(eventBridgeInput)

	if err != nil {
		return err
	}

	if *res.FailedEntryCount > 0 {
		return fmt.Errorf("failed to put events: %s", *res.Entries[0].ErrorMessage)
	}

	return nil
}
