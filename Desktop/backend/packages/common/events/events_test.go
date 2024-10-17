package events

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventsPublisherSuite struct {
	suite.Suite
}

func (s *EventsPublisherSuite) TestCallsEventBridgeProperly() {
	eventbridgeApi := new(mocks.EventBridgeAPI)
	publisher := EventPublisher{
		EventBridge: eventbridgeApi,
	}

	eventbridgeApi.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{FailedEntryCount: aws.Int64(0)}, nil)

	err := publisher.PublishConnectionCreatedEvent("test", "test-service")

	eventbridgeApi.AssertCalled(s.T(), "PutEvents", mock.MatchedBy(func(input *eventbridge.PutEventsInput) bool {
		entry := input.Entries[0]
		detailBytes := []byte(*entry.Detail)
		var event ConnectionCreatedEvent
		_ = json.Unmarshal(detailBytes, &event)

		return *entry.DetailType == eventconstants.CONNECTION_CREATED_EVENT &&
			*entry.Source == "test-service" &&
			event.UserID == "test"
	}))

	s.Nil(err)
}

func (s *EventsPublisherSuite) TestFailureCountIsHandled() {
	eventbridgeApi := new(mocks.EventBridgeAPI)
	publisher := EventPublisher{
		EventBridge: eventbridgeApi,
	}

	eventbridgeApi.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: aws.Int64(1),
		Entries: []*eventbridge.PutEventsResultEntry{
			{ErrorMessage: aws.String("uh-oh")},
		},
	}, nil)

	err := publisher.PublishConnectionCreatedEvent("test", "test-service")

	s.NotNil(err)
}

func (s *EventsPublisherSuite) TestErrorHandled() {
	eventbridgeApi := new(mocks.EventBridgeAPI)
	publisher := EventPublisher{
		EventBridge: eventbridgeApi,
	}

	eventbridgeApi.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{FailedEntryCount: aws.Int64(1)}, fmt.Errorf("lol"))

	err := publisher.PublishConnectionCreatedEvent("test", "test-service")

	s.NotNil(err)
}

func TestEventPublisher(t *testing.T) {
	suite.Run(t, new(EventsPublisherSuite))
}
