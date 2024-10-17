package main

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CountFormsSentTestSuite struct {
	suite.Suite
	gormDb *gorm.DB
	mockDb sqlmock.Sqlmock
}

func (suite *CountFormsSentTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *CountFormsSentTestSuite) TestCountFormsSent_Success() {
	const sqlCount = `SELECT count(*) FROM "form_sents"`

	suite.mockDb.ExpectQuery(sqlCount).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(2))

	actual, err := countFormsSent(suite.gormDb)

	suite.NotNil(actual)
	suite.Nil(err)
	suite.EqualValues(actual, 2)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *CountFormsSentTestSuite) TestCountFormsSent_Fail() {
	const sqlCount = `SELECT count(*) FROM "form_sents"`

	suite.mockDb.ExpectQuery(sqlCount).WillReturnError(fmt.Errorf("some error"))

	_, err := countFormsSent(suite.gormDb)

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

/* Execute Suites */

func TestCountFormsSentTestSuite(t *testing.T) {
	suite.Run(t, new(CountFormsSentTestSuite))
}
