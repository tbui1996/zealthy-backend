package main

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ChatSessionUpdateStatusSuite struct {
	suite.Suite
	Request      UpdateStatusRequest
	Insert       string
	SQLMock      sqlmock.Sqlmock
	ClosedUpdate string
	OpenUpdate   string
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (suite *ChatSessionUpdateStatusSuite) SetupTest() {
	db, theMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.SQLMock = theMock

	suite.Insert = `INSERT INTO "session_descriptors" ("session_id","name","value") VALUES ($1,$2,$3) ON CONFLICT ("session_id","name") DO UPDATE SET "value"="excluded"."value"`
	suite.ClosedUpdate = `UPDATE "session_statuses" SET "closed_at"=$1,"status"=$2 WHERE session_id = $3`
	suite.OpenUpdate = `UPDATE "session_statuses" SET "opened_at"=$1,"status"=$2 WHERE session_id = $3`

	suite.Request = UpdateStatusRequest{
		DB:            gdb,
		Logger:        zaptest.NewLogger(suite.T()),
		Open:          false,
		RideScheduled: false,
		SessionId:     "1",
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusSuccessClosed() {
	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.RIDE_SCHEDULED, "false").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.SQLMock.ExpectExec(suite.ClosedUpdate).
		WithArgs(AnyTime{}, model.CLOSED, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.SQLMock.ExpectCommit()

	err := Handler(suite.Request)

	suite.NoError(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusSuccessOpen() {
	req := suite.Request
	req.Open = true

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.RIDE_SCHEDULED, "false").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.SQLMock.ExpectExec(suite.OpenUpdate).
		WithArgs(AnyTime{}, model.OPEN, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.SQLMock.ExpectCommit()

	err := Handler(req)

	suite.NoError(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusSuccessRideScheduled() {
	req := suite.Request
	req.RideScheduled = true

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.RIDE_SCHEDULED, "true").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.SQLMock.ExpectExec(suite.ClosedUpdate).
		WithArgs(AnyTime{}, model.CLOSED, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.SQLMock.ExpectCommit()

	err := Handler(req)

	suite.NoError(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusError() {
	req := suite.Request
	req.RideScheduled = true

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.RIDE_SCHEDULED, "true").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.SQLMock.ExpectExec(suite.ClosedUpdate).
		WithArgs(AnyTime{}, model.CLOSED, 1).
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.SQLMock.ExpectRollback()

	err := Handler(req)

	suite.Error(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusErrorRideScheduled() {
	req := suite.Request
	req.RideScheduled = true

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.RIDE_SCHEDULED, "true").
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.SQLMock.ExpectRollback()

	err := Handler(req)

	suite.Error(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionUpdateStatusSuite) Test_UpdateStatusTransactionError() {
	req := suite.Request

	suite.SQLMock.ExpectBegin().WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(req)

	suite.Error(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatSessionsGet(t *testing.T) {
	suite.Run(t, new(ChatSessionUpdateStatusSuite))
}
