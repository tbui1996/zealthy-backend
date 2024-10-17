package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	commonEvents "github.com/circulohealth/sonar-backend/packages/common/events"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/pretty"
)

func connect(context context.Context, event events.CloudWatchEvent) (commonEvents.MessageSentEvent, error) {
	var details commonEvents.MessageSentEvent

	logger, err := logging.NewLoggerFromEvent(event)
	if err != nil {
		return details, err
	}
	defer logging.SyncLogger(logger)

	logger.Debug(pretty.Sprint(event))

	if event.DetailType != eventconstants.MESSAGE_SENT_EVENT {
		errorMessage := fmt.Sprintf("expected detail type: message_sent_event, got: %s", event.DetailType)
		logger.Error(errorMessage)
		return details, errors.New(errorMessage)
	}

	err = json.Unmarshal(event.Detail, &details)

	if err != nil {
		errorMessage := fmt.Sprintf("failed to parse event details: %s", err)
		logger.Error(errorMessage)
		return details, errors.New(errorMessage)
	}

	if details.ReceiverId == "" {
		errorMessage := fmt.Sprintf("invalid event details, ReceiverId is required. %s", string(event.Detail))
		logger.Error(errorMessage)
		return details, errors.New(errorMessage)
	}

	return details, nil
}

func main() {
	lambda.Start(connect)
}
