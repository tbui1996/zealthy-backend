package main

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UpdateChatNotesTestSuite struct {
	suite.Suite
	db      *gorm.DB
	theMock sqlmock.Sqlmock
}

func (suite *UpdateChatNotesTestSuite) SetupTest() {
	db, theMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = gormDB
	suite.theMock = theMock
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotesUpdate_Success() {
	row := sqlmock.
		NewRows([]string{"session_id", "notes"}).
		AddRow(1, "hello world")

	suite.theMock.ExpectBegin()
	suite.theMock.ExpectQuery(`SELECT * FROM "session_notes" WHERE "session_notes"."session_id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(row)

	suite.theMock.ExpectExec(`UPDATE "session_notes" SET "notes"=$1 WHERE session_id = $2`).
		WithArgs("goodbye world", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.theMock.ExpectCommit()

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "goodbye world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotesCreate_Success() {
	suite.theMock.ExpectBegin()
	suite.theMock.ExpectQuery(`SELECT * FROM "session_notes" WHERE "session_notes"."session_id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	suite.theMock.ExpectExec(`INSERT INTO "session_notes" ("session_id","notes") VALUES ($1,$2)`).
		WithArgs(1, "hello world").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.theMock.ExpectCommit()

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "hello world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotesTake_Fail() {
	suite.theMock.ExpectBegin()
	suite.theMock.ExpectQuery(`SELECT * FROM "session_notes" WHERE "session_notes"."session_id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.theMock.ExpectRollback()

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "hello world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotesCreate_Fail() {
	suite.theMock.ExpectBegin()
	suite.theMock.ExpectQuery(`SELECT * FROM "session_notes" WHERE "session_notes"."session_id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	suite.theMock.ExpectExec(`INSERT INTO "session_notes" ("session_id","notes") VALUES ($1,$2)`).
		WithArgs(1, "hello world").
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.theMock.ExpectRollback()

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "hello world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotesUpdate_Fail() {
	row := sqlmock.
		NewRows([]string{"session_id", "notes"}).
		AddRow(1, "hello world")

	suite.theMock.ExpectBegin()
	suite.theMock.ExpectQuery(`SELECT * FROM "session_notes" WHERE "session_notes"."session_id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(row)

	suite.theMock.ExpectExec(`UPDATE "session_notes" SET "notes"=$1 WHERE session_id = $2`).
		WithArgs("goodbye world", 1).
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.theMock.ExpectRollback()

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "goodbye world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UpdateChatNotesTestSuite) TestUpdateChatNotes_FailOnBegin() {
	suite.theMock.ExpectBegin().WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(UpdateChatNotesRequest{
		SessionID: 1,
		Notes:     "goodbye world",
		DB:        suite.db,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)

	if err := suite.theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

/* Execute Suites */

func TestUpdateChatNotesTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateChatNotesTestSuite))
}
