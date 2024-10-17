package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
}

func (suite *HandlerTestSuite) TestGetForm_ShouldReturnErrorIfFormDoesntExist() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	date := time.Now()
	id := "1"
	const sqlUpdate = `UPDATE "forms" SET "date_closed"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, id).WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err := Handler(gdb, date, id)

	suite.NotNil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *HandlerTestSuite) TestSaveDateClosed_ShouldWriteToDateClosedColumn() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	date := time.Now()
	id := "1"
	const sqlUpdate = `UPDATE "forms" SET "date_closed"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, id).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := Handler(gdb, date, id)

	suite.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
