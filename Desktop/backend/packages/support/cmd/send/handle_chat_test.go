package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandleChatTestSuite struct {
	suite.Suite
}

func upsertFromConstant() string {
	upsert := model.LastMessageUpsert
	upsert = strings.Replace(upsert, "@message", "$1", 1)
	upsert = strings.Replace(upsert, "@sent", "$2", 1)
	upsert = strings.Replace(upsert, "@id", "$3", 1)
	upsert = strings.Replace(upsert, "@user", "$4", 1)
	return upsert
}

func (suite *HandleChatTestSuite) TestHandleChat_SendsMessage() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	eventbridgeMock := new(mocks.EventBridgeAPI)
	eventbridgeMock.On("PutEvents", mock.Anything).Return(&eventbridge.PutEventsOutput{
		FailedEntryCount: aws.Int64(0),
	}, nil)

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	m := request.Chat{Session: "1", Sender: "Okta_test", Message: "Test", File: ""}

	expected := &model.ChatMessage{ID: "1", SessionID: m.Session, SenderID: m.Sender, Message: m.Message, CreatedTimestamp: 1, FileID: nil}

	mockChatMessageRepository := new(session.MockChatMessageRepository)
	mockChatMessageRepository.On("Create", mock.Anything).Return(expected, nil)

	repo := &session.RDSChatSessionRepository{DB: theDB, ChatMessageRepository: mockChatMessageRepository, EventBridge: eventbridgeMock}

	row := sqlmock.NewRows([]string{"id", "user_id", "session_id"}).AddRow(1, "Okta_test", 1).AddRow(2, "test", 1)

	theMock.ExpectBegin()

	theMock.ExpectQuery("SELECT * FROM \"session_users\" WHERE session_id = $1").WithArgs("1").WillReturnRows(row)

	sRow := sqlmock.NewRows([]string{"session_id", "status"}).AddRow(1, model.OPEN)

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 LIMIT 1").WithArgs("1").WillReturnRows(sRow)

	rRow := sqlmock.NewRows([]string{"id", "created", "chat_type"}).AddRow(1, time.Now(), model.CIRCULATOR)

	theMock.ExpectQuery("SELECT * FROM \"sessions\" WHERE id = $1 LIMIT 1").WithArgs("1").WillReturnRows(rRow)

	theMock.ExpectCommit()

	theMock.ExpectBegin()
	theMock.ExpectExec(upsertFromConstant()).
		WithArgs("Test", 1, "1", "Okta_test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectCommit()

	config := requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zap.NewExample(),
		Event: events.APIGatewayWebsocketProxyRequest{
			RequestContext: events.APIGatewayWebsocketProxyRequestContext{
				Authorizer: map[string]interface{}{}},
		}}

	message, _ := json.Marshal(m)

	client := new(router.MockRouter)
	client.On("Send", mock.Anything).Return(nil)

	mockClient := &router.Session{
		Router: client,
	}

	resErr := HandleChat(&config, string(message), repo, mockClient)

	suite.Nil(resErr)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	mockChatMessageRepository.AssertCalled(suite.T(), "Create", m)
}

func TestHandleChatTestSuite(t *testing.T) {
	suite.Run(t, new(HandleChatTestSuite))
}
