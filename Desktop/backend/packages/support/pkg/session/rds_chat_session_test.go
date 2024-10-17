package session

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/circulohealth/sonar-backend/packages/common/events"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RDSChatSessionTestSuite struct {
	suite.Suite
}

const CIRCULO = "circulo"

func upsertFromConstant() string {
	upsert := model.LastMessageUpsert
	upsert = strings.Replace(upsert, "@message", "$1", 1)
	upsert = strings.Replace(upsert, "@sent", "$2", 1)
	upsert = strings.Replace(upsert, "@id", "$3", 1)
	upsert = strings.Replace(upsert, "@user", "$4", 1)
	return upsert
}

func CreateBasicMockEventBridge() *mocks.EventBridgeAPI {
	eventbridgeMock := new(mocks.EventBridgeAPI)
	var failedCount int64 = 0
	eventbridgeMock.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: &failedCount,
	}, nil)
	return eventbridgeMock
}

func (suite *RDSChatSessionTestSuite) TestPostRecordLatest_ShouldInsertLastMessage() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	eventbridgeMock := CreateBasicMockEventBridge()
	internalUserID := "2"

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	sess := &RDSChatSession{
		db: gormDB,
		dto: &model.ChatSessionDTO{
			ID:             "1",
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := &model.ChatMessage{
		CreatedTimestamp: 1,
		Message:          "test",
		SenderID:         CIRCULO,
	}

	mock.ExpectBegin()
	mock.ExpectExec(upsertFromConstant()).
		WithArgs("test", 1, "1", CIRCULO).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updated, err := sess.recordLastMessage(message)
	suite.Equal(int64(1), *updated)
	suite.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionTestSuite) TestPostRecordLatest_ShouldUpdateLastMessage() {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	eventbridgeMock := CreateBasicMockEventBridge()

	internalUserID := "2"

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	sess := &RDSChatSession{
		db: gormDB,
		dto: &model.ChatSessionDTO{
			ID:             "1",
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := &model.ChatMessage{
		CreatedTimestamp: 2,
		Message:          "update",
		SenderID:         CIRCULO,
	}

	sqlmock.NewRows([]string{"session_user_id", "last_message", "last_sent"}).
		AddRow(1, "test", 1)

	dbMock.ExpectBegin()
	dbMock.ExpectExec(upsertFromConstant()).
		WithArgs("update", 2, "1", CIRCULO).
		WillReturnResult(sqlmock.NewResult(0, 1))
	dbMock.ExpectCommit()

	updated, err := sess.recordLastMessage(message)
	suite.Equal(int64(1), *updated)
	suite.Nil(err)

	if err := dbMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionTestSuite) TestPostRecordLatest_ShouldPublishLoopSenderType() {
	eventbridgeMock := new(mocks.EventBridgeAPI)
	eventbridgeMock.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: aws.Int64(0),
	}, nil)

	internalUserID := "2"

	sess := &RDSChatSession{
		dto: &model.ChatSessionDTO{
			ID:             "1",
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := &model.ChatMessage{
		CreatedTimestamp: 2,
		Message:          "update",
		SenderID:         CIRCULO,
	}

	err := sess.publishMessageSentEvent(message)
	suite.Nil(err)

	eventbridgeMock.AssertCalled(suite.T(), "PutEvents", mock.MatchedBy(func(input *eventbridge.PutEventsInput) bool {
		detailsBytes := []byte(*input.Entries[0].Detail)
		var details events.MessageSentEvent = events.MessageSentEvent{}

		_ = json.Unmarshal(detailsBytes, &details)
		return *input.Entries[0].DetailType == eventconstants.MESSAGE_SENT_EVENT && details.SenderType == events.LoopSenderType
	}))
}

func (suite *RDSChatSessionTestSuite) TestPostRecordLatest_ShouldPublishInternalSenderType() {
	eventbridgeMock := new(mocks.EventBridgeAPI)
	eventbridgeMock.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: aws.Int64(0),
	}, nil)

	internalUserID := CIRCULO // when this equals the message sender id the sender type is internal

	sess := &RDSChatSession{
		dto: &model.ChatSessionDTO{
			ID:             "1",
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := &model.ChatMessage{
		CreatedTimestamp: 2,
		Message:          "update",
		SenderID:         CIRCULO,
	}

	err := sess.publishMessageSentEvent(message)
	suite.Nil(err)

	eventbridgeMock.AssertCalled(suite.T(), "PutEvents", mock.MatchedBy(func(input *eventbridge.PutEventsInput) bool {
		detailsBytes := []byte(*input.Entries[0].Detail)
		var details events.MessageSentEvent = events.MessageSentEvent{}

		_ = json.Unmarshal(detailsBytes, &details)
		return *input.Entries[0].DetailType == eventconstants.MESSAGE_SENT_EVENT && details.SenderType == events.InternalSenderType
	}))
}

func (suite *RDSChatSessionTestSuite) TestPostRecordLatest_ShouldHandleFailedEntryCount() {
	eventbridgeMock := new(mocks.EventBridgeAPI)
	eventbridgeMock.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: aws.Int64(1),
	}, nil)

	internalUserID := "2"

	sess := &RDSChatSession{
		dto: &model.ChatSessionDTO{
			ID:             "1",
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := &model.ChatMessage{
		CreatedTimestamp: 2,
		Message:          "update",
		SenderID:         CIRCULO,
	}

	err := sess.publishMessageSentEvent(message)
	suite.NotNil(err)
}

func (suite *RDSChatSessionTestSuite) TestIsValidSender_ShouldReturnTrueIfMessageSenderIsUser() {
	userID := "1"
	internalUserID := "2"
	sess := &RDSChatSession{
		dto: &model.ChatSessionDTO{
			ID:             "1",
			UserID:         userID,
			InternalUserID: &internalUserID,
		},
	}

	message := request.Chat{
		Message: "test",
		Sender:  userID,
		Session: "1",
		File:    "",
	}

	valid := sess.isValidSender(message)

	suite.True(valid)
}

func (suite *RDSChatSessionTestSuite) TestIsValidSender_ShouldReturnTrueIfMessageSenderIsInternalUser() {
	userID := "1"
	internalUserID := "2"
	sess := &RDSChatSession{
		pending: false,
		dto: &model.ChatSessionDTO{
			ID:             "1",
			UserID:         userID,
			InternalUserID: &internalUserID,
		},
	}

	message := request.Chat{
		Message: "test",
		Sender:  internalUserID,
		Session: "1",
		File:    "",
	}

	valid := sess.isValidSender(message)

	suite.True(valid)
}

func (suite *RDSChatSessionTestSuite) TestAppendRequestMessage_ShouldErrorIfSenderIsNotInSession() {
	eventbridgeMock := CreateBasicMockEventBridge()
	userID := "1"
	internalUserID := "2"
	sess := &RDSChatSession{
		dto: &model.ChatSessionDTO{
			ID:             "1",
			UserID:         userID,
			InternalUserID: &internalUserID,
		},
		eventBridge: eventbridgeMock,
	}

	message := request.Chat{
		Message: "test",
		Sender:  "3",
		Session: "1",
		File:    "",
	}

	actual, err := sess.AppendRequestMessage(message)

	suite.Nil(actual)
	suite.Error(err)
}

func (suite *RDSChatSessionTestSuite) TestAppendRequestMessage_ShouldCallCreateInChatMessageRepository() {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	eventbridgeMock := CreateBasicMockEventBridge()

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	userID := CIRCULO
	message := request.Chat{
		Message: "test",
		Sender:  userID,
		Session: "1",
		File:    "",
	}

	expected := &model.ChatMessage{
		ID:               "1",
		SessionID:        message.Session,
		SenderID:         message.Sender,
		Message:          message.Message,
		CreatedTimestamp: 1,
		FileID:           nil,
	}

	mockChatMessageRepository := new(MockChatMessageRepository)
	mockChatMessageRepository.On("Create", mock.Anything).Return(expected, nil)

	internalUserID := "2"
	sess := &RDSChatSession{
		db: gormDB,
		dto: &model.ChatSessionDTO{
			ID:             "1",
			UserID:         userID,
			InternalUserID: &internalUserID,
		},
		chatMessageRepository: mockChatMessageRepository,
		eventBridge:           eventbridgeMock,
	}

	dbMock.ExpectBegin()
	dbMock.ExpectExec(upsertFromConstant()).
		WithArgs("test", 1, "1", CIRCULO).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	actual, err := sess.AppendRequestMessage(message)

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Equal(expected, actual)
	mockChatMessageRepository.AssertCalled(suite.T(), "Create", message)

	if err := dbMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRDSChatSessionTestSuite(t *testing.T) {
	suite.Run(t, new(RDSChatSessionTestSuite))
}
