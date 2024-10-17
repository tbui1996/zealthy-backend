package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DiscardAndSubmitTestSuite struct {
	suite.Suite
	sqlQueryDiscard     string
	sqlQueryInputSubmit string
	gormDb              *gorm.DB
	mockDb              sqlmock.Sqlmock
	formDiscard         model.FormDiscard
	formSubArray        []model.FormSubmission
	inputSub            model.InputSubmission
}

func (suite *DiscardAndSubmitTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *DiscardAndSubmitTestSuite) SetupSuite() {
	suite.sqlQueryDiscard = `SELECT * FROM "form_discards" WHERE form_sent_id IN ($1)`
	suite.sqlQueryInputSubmit = `SELECT * FROM "input_submissions" WHERE form_submission_id = $1`
	suite.formDiscard = model.FormDiscard{
		ID:         1,
		FormSentId: 1,
		Deleted:    time.Now(),
	}
	suite.formSubArray = []model.FormSubmission{
		{
			ID:         1,
			FormSentId: 1,
		},
	}
	suite.inputSub = model.InputSubmission{
		ID:               1,
		FormSubmissionId: 1,
		InputId:          1,
		Response:         "test",
	}
}
func (suite *DiscardAndSubmitTestSuite) TestDiscardAndSubmit_Success() {
	discardRows := sqlmock.NewRows([]string{"id", "form_sent_id", "deleted"}).AddRow(suite.formDiscard.ID, suite.formDiscard.FormSentId, suite.formDiscard.Deleted)
	inputSubmitRows := sqlmock.NewRows([]string{"id", "form_submission_id", "input_id", "response"}).AddRow(suite.inputSub.ID, suite.inputSub.FormSubmissionId, suite.inputSub.InputId, suite.inputSub.Response)
	ids := []int{1}
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlQueryDiscard).WithArgs(1).WillReturnRows(discardRows)
	suite.mockDb.ExpectQuery(suite.sqlQueryInputSubmit).WithArgs(1).WillReturnRows(inputSubmitRows)
	suite.mockDb.ExpectCommit()

	actualDiscards, actualInputs, err := getDiscardAndSubmitValues(&DiscardAndSubmitInput{
		Db:        suite.gormDb,
		FormSents: ids,
		Submit:    suite.formSubArray,
	})

	suite.Nil(err)
	suite.Equal(actualDiscards, []model.FormDiscard{suite.formDiscard})
	suite.Equal(actualInputs, [][]model.InputSubmission{{suite.inputSub}})

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *DiscardAndSubmitTestSuite) TestDiscardAndSubmit_ErrorOnDiscard() {
	ids := []int{1}
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlQueryDiscard).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	actualDiscards, actualInputs, err := getDiscardAndSubmitValues(&DiscardAndSubmitInput{
		Db:        suite.gormDb,
		FormSents: ids,
		Submit:    suite.formSubArray,
	})

	suite.NotNil(err)
	suite.Nil(actualDiscards)
	suite.Nil(actualInputs)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *DiscardAndSubmitTestSuite) TestDiscardAndSubmit_ErrorOnInputs() {
	discardRows := sqlmock.NewRows([]string{"id", "form_sent_id", "deleted"}).AddRow(suite.formDiscard.ID, suite.formDiscard.FormSentId, suite.formDiscard.Deleted)

	ids := []int{1}
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlQueryDiscard).WithArgs(1).WithArgs(1).WillReturnRows(discardRows)
	suite.mockDb.ExpectQuery(suite.sqlQueryInputSubmit).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	actualDiscards, actualInputs, err := getDiscardAndSubmitValues(&DiscardAndSubmitInput{
		Db:        suite.gormDb,
		FormSents: ids,
		Submit:    suite.formSubArray,
	})

	suite.NotNil(err)
	suite.Nil(actualDiscards)
	suite.Nil(actualInputs)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

type FindFormSentTestSuite struct {
	suite.Suite
	sqlQuery string
	gormDb   *gorm.DB
	mockDb   sqlmock.Sqlmock
	sent     model.FormSent
}

func (suite *FindFormSentTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *FindFormSentTestSuite) SetupSuite() {
	suite.sqlQuery = `SELECT * FROM "form_sents" WHERE form_id = $1`
	suite.sent = model.FormSent{
		ID:     1,
		FormId: 1,
		Sent:   time.Now(),
	}
}

func (suite *FindFormSentTestSuite) TestFindFormSent_Success() {
	rows := sqlmock.NewRows([]string{"id", "form_id", "sent"}).AddRow(suite.sent.ID, suite.sent.FormId, suite.sent.Sent)

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnRows(rows)
	const formID = "1"
	actual, err := findFormSent(formID, suite.gormDb)

	suite.Nil(err)
	suite.NotNil(actual)
	suite.Equal(actual, []model.FormSent{suite.sent})

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *FindFormSentTestSuite) TestFindFormSent_Fail() {
	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnError(fmt.Errorf("some error"))
	const formID = "1"
	_, err := findFormSent(formID, suite.gormDb)

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

type FindSubmitByFormSentTestSuite struct {
	suite.Suite
	sqlQuery   string
	gormDb     *gorm.DB
	mockDb     sqlmock.Sqlmock
	submission model.FormSubmission
}

func (suite *FindSubmitByFormSentTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *FindSubmitByFormSentTestSuite) SetupSuite() {
	suite.sqlQuery = `SELECT * FROM "form_submissions" WHERE form_sent_id IN ($1)`
	suite.submission = model.FormSubmission{
		ID:         1,
		FormSentId: 1,
	}
}

func (suite *FindSubmitByFormSentTestSuite) TestFindSubmitByFormSent_Success() {
	rows := sqlmock.NewRows([]string{"id", "form_sent_id"}).AddRow(suite.submission.ID, suite.submission.FormSentId)

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnRows(rows)
	var formSents = []int{1}
	actual, err := findSubmitByFormSent(formSents, suite.gormDb)

	suite.Nil(err)
	suite.NotNil(actual)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *FindSubmitByFormSentTestSuite) TestFindSubmitByFormSent_Fail() {
	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnError(fmt.Errorf("some error"))
	var formSents = []int{1}
	_, err := findSubmitByFormSent(formSents, suite.gormDb)

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

/* Execute Suites */

func TestDiscardAndSubmitTestSuite(t *testing.T) {
	suite.Run(t, new(DiscardAndSubmitTestSuite))
}

func TestFindFormSentTestSuite(t *testing.T) {
	suite.Run(t, new(FindFormSentTestSuite))
}

func TestFindSubmitByFormSentTestSuite(t *testing.T) {
	suite.Run(t, new(FindSubmitByFormSentTestSuite))
}
