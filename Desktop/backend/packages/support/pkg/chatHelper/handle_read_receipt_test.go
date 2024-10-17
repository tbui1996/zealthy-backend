package chatHelper

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
	"time"
)

type HandleReadReceiptTestSuite struct {
	suite.Suite
}

func (suite *HandleReadReceiptTestSuite) TestHandleReadReceipt_ChatOpen() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &session.RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(session.MockChatMessageRepository),
	}

	row := sqlmock.NewRows([]string{"session_id", "status"}).
		AddRow(1, "OPEN")

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 ORDER BY \"session_statuses\".\"session_id\" LIMIT 1").
		WithArgs("1").
		WillReturnRows(row)

	read := time.Now().Unix()

	query := model.LastReadUpsert
	query = strings.Replace(query, "@read", "$1", 1)
	query = strings.Replace(query, "@id", "$2", 1)
	query = strings.Replace(query, "@user", "$3", 1)

	theMock.ExpectBegin()
	theMock.ExpectExec(query).
		WithArgs(read, "1", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectCommit()

	req := request.ReadReceiptRequest{SessionID: "1", UserID: "test"}

	b, _ := json.Marshal(req)

	err = HandleReadReceipt(string(b), read, repo)

	suite.Nil(err)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *HandleReadReceiptTestSuite) TestHandleReadReceipt_ChatPending() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &session.RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(session.MockChatMessageRepository),
	}

	row := sqlmock.NewRows([]string{"session_id", "status"}).
		AddRow(1, "PENDING")

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 ORDER BY \"session_statuses\".\"session_id\" LIMIT 1").
		WithArgs("1").
		WillReturnRows(row)

	read := time.Now().Unix()

	query := model.LastReadUpsert
	query = strings.Replace(query, "@read", "$1", 1)
	query = strings.Replace(query, "@id", "$2", 1)
	query = strings.Replace(query, "@user", "$3", 1)

	theMock.ExpectBegin()
	theMock.ExpectExec(query).
		WithArgs(read, "1", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectCommit()

	req := request.ReadReceiptRequest{SessionID: "1", UserID: "test"}

	b, _ := json.Marshal(req)

	err = HandleReadReceipt(string(b), read, repo)

	suite.Nil(err)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleReadReceiptTestSuite(t *testing.T) {
	suite.Run(t, new(HandleReadReceiptTestSuite))
}
