package session

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/circulohealth/sonar-backend/packages/common/events"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"gorm.io/gorm"
)

type RDSChatSession struct {
	db                    *gorm.DB
	dto                   *model.ChatSessionDTO
	chatMessageRepository iface.ChatMessageRepository
	dbChatMessages        *dynamodb.DynamoDB
	eventBridge           eventbridgeiface.EventBridgeAPI

	pending bool
}

func (sess *RDSChatSession) ToResponseDTO() response.ChatSessionResponseDTO {
	return response.ChatSessionResponseDTO{
		IsPending:            sess.pending,
		ID:                   sess.dto.ID,
		UserID:               sess.dto.UserID,
		InternalUserID:       sess.InternalUserID(),
		Topic:                sess.dto.Topic,
		Patient:              sess.Patient(),
		ChatType:             model.ChatTypeToString(sess.dto.ChatType),
		Notes:                sess.dto.Notes,
		CreatedTimestamp:     sess.dto.CreatedTimestamp,
		UserLastRead:         sess.UserLastRead(),
		InternalUserLastRead: sess.InternalUserLastRead(),
		ChatOpen:             sess.IsOpen(),
		LastMessageTimestamp: sess.dto.LastMessageTimestamp,
		LastMessagePreview:   sess.dto.LastMessagePreview,
		LastMessageSenderID:  sess.dto.LastMessageSenderID,
		Starred:              sess.dto.Starred,
	}
}

func (sess *RDSChatSession) ID() string {
	return sess.dto.ID
}

func (sess *RDSChatSession) UserID() string {
	return sess.dto.UserID
}

func (sess *RDSChatSession) LastMessageTimestamp() int64 {
	return sess.dto.LastMessageTimestamp
}

func (sess *RDSChatSession) UserLastRead() int64 {
	return sess.dto.UserLastRead
}

func (sess *RDSChatSession) InternalUserLastRead() int64 {
	if sess.pending {
		return 0
	}

	return sess.dto.InternalUserLastRead
}

func (sess *RDSChatSession) InternalUserID() string {
	if sess.pending {
		return ""
	}

	if sess.dto.InternalUserID == nil {
		return ""
	}

	return *sess.dto.InternalUserID
}

func (sess *RDSChatSession) Patient() response.PatientResponse {
	return response.PatientResponse{
		ID:          sess.dto.Patient.ID,
		Name:        sess.dto.Patient.Name,
		LastName:    sess.dto.Patient.LastName,
		Address:     sess.dto.Patient.Address,
		InsuranceID: sess.dto.Patient.InsuranceID,
		Birthday:    sess.dto.Patient.Birthday,
	}
}

func (sess *RDSChatSession) IsPending() bool {
	return sess.pending
}

func (sess *RDSChatSession) IsOpen() bool {
	if sess.pending {
		return true
	}

	return sess.dto.ChatOpen
}
func (sess *RDSChatSession) GetMessages() ([]model.ChatMessage, error) {
	return sess.chatMessageRepository.GetMessagesForSession(sess.ID())
}

func (sess *RDSChatSession) Type() model.ChatType {
	return sess.dto.ChatType
}

func (sess *RDSChatSession) AppendRequestMessage(message request.Chat) (*model.ChatMessage, error) {
	// validate the sender is a chat session member
	if !sess.isValidSender(message) {
		return nil, fmt.Errorf("only session members can send messages in the session. Sender: [%s], User: [%s], InternalUser: [%s]", message.Sender, sess.UserID(), sess.InternalUserID())
	}

	created, err := sess.chatMessageRepository.Create(message)

	if err != nil {
		return nil, err
	}

	_, err = sess.recordLastMessage(created)
	if err != nil {
		return nil, err
	}

	err = sess.publishMessageSentEvent(created)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (sess *RDSChatSession) isValidSender(message request.Chat) bool {
	return message.Sender == sess.UserID() || message.Sender == sess.InternalUserID()
}

func (sess *RDSChatSession) recordLastMessage(message *model.ChatMessage) (*int64, error) {
	db := sess.db

	tx := db.Begin()
	result := tx.Exec(
		model.LastMessageUpsert,
		sql.Named("message", message.Message),
		sql.Named("sent", message.CreatedTimestamp),
		sql.Named("id", sess.ID()),
		sql.Named("user", message.SenderID),
	)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	tx.Commit()

	return &result.RowsAffected, nil
}

func (sess *RDSChatSession) publishMessageSentEvent(message *model.ChatMessage) error {

	// TODO: when the use case arises we will need to send events when InternalUserID is null and handle them in the offline-message-notifier step function
	if sess.dto.InternalUserID != nil {
		receiverId, senderId, senderType := *sess.dto.InternalUserID, message.SenderID, events.LoopSenderType
		if receiverId == senderId {
			receiverId = sess.dto.UserID
			senderType = events.InternalSenderType
		}

		event := events.MessageSentEvent{
			SenderId:   senderId,
			ReceiverId: receiverId,
			SenderType: senderType,
			SentAt:     message.CreatedTimestamp,
			SessionId:  sess.dto.ID,
			MessageId:  message.ID,
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
					DetailType:   aws.String(eventconstants.MESSAGE_SENT_EVENT),
					Resources:    []*string{},
					Source:       aws.String(eventconstants.SUPPORT_SERVICE),
				},
			},
		}

		output, err := sess.eventBridge.PutEvents(eventBridgeInput)
		if err != nil {
			return err
		}

		if *output.FailedEntryCount > 0 {
			return fmt.Errorf("failed to post to event bridge: %s", output.String())
		}
	}
	return nil
}
