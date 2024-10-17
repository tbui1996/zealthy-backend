package main

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type ChatSessionStarSuite struct {
	suite.Suite
	Request SubmitChatSessionStarRequest
	SQLMock sqlmock.Sqlmock
	Insert  string
}

func (suite *ChatSessionStarSuite) SetupTest() {
	db, theMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.SQLMock = theMock

	suite.Insert = `INSERT INTO "session_descriptors" ("session_id","name","value") VALUES ($1,$2,$3) ON CONFLICT ("session_id","name") DO UPDATE SET "value"="excluded"."value"`

	suite.Request = SubmitChatSessionStarRequest{
		DB:        gdb,
		SessionID: "1",
		OnStar:    true,
		Logger:    zaptest.NewLogger(suite.T()),
	}
}

func (suite *ChatSessionStarSuite) TestOnStar_Success() {
	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.STARRED, "true").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.SQLMock.ExpectCommit()

	err := Handler(suite.Request)

	suite.Nil(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionStarSuite) TestOnStar_Error() {
	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectExec(suite.Insert).
		WithArgs(1, model.STARRED, "true").
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))
	suite.SQLMock.ExpectRollback()

	err := Handler(suite.Request)

	suite.NotNil(err)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *ChatSessionStarSuite) TestOnStar_ErrorNotInt() {
	req := suite.Request
	req.SessionID = "test"
	err := Handler(req)

	suite.NotNil(err)
}

func TestStarChat(t *testing.T) {
	suite.Run(t, new(ChatSessionStarSuite))
}
