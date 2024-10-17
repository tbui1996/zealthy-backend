package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/common/pretty"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"go.uber.org/zap"
)

func Handler(message events.SQSMessage, repo iface.ChatSessionRepository) error {
	var lf logging.LoggerFields
	logger := logging.Must(lf.FromSQSMessage(message))

	connectionID, ok := message.MessageAttributes["ConnectionId"]
	if !ok {
		logger.Error("couldn't get connection id")
	} else {
		logger = logger.With(zap.String("connectionID", *connectionID.StringValue), zap.String("messageID", message.MessageId))
		logger.Info("processing support receive message from " + *connectionID.StringValue + ": " + message.MessageId)
	}
	var supportRequest request.SupportRequestReceive
	if err := json.Unmarshal([]byte(message.Body), &supportRequest); err != nil {
		return fmt.Errorf("unable to unmarshal body: %+v (%s)", message.Body, err)
	}

	switch supportRequest.Type {
	case "read_receipt":
		err := chatHelper.HandleReadReceipt(supportRequest.Message, time.Now().Unix(), repo)
		if err != nil {
			return err
		}
	default:
		logger.Info("unsupported request type in support request", zap.String("supportRequest", pretty.Sprint(supportRequest)))
	}

	return nil
}
