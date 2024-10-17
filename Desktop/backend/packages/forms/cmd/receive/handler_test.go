package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type discardTestSuite struct {
	suite.Suite
	sqlInsert      string
	gormDb         *gorm.DB
	mockDb         sqlmock.Sqlmock
	formDiscard    model.FormDiscard
	discardRequest request.DiscardFormRequest
}

func (suite *discardTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *discardTestSuite) SetupSuite() {
	const formSentId = 1
	suite.sqlInsert = `INSERT INTO "form_discards" ("form_sent_id","deleted") VALUES ($1,$2) RETURNING "id"`
	suite.formDiscard = model.FormDiscard{
		FormSentId: formSentId,
		Deleted:    time.Now(),
	}
	suite.discardRequest = request.DiscardFormRequest{
		FormSentId: formSentId,
	}
}

func (suite *discardTestSuite) TestDiscard_Success() {
	var e error
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.formDiscard.FormSentId, suite.formDiscard.Deleted).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	err := discard(&DiscardInput{
		Db:             suite.gormDb,
		Err:            e,
		DiscardRequest: suite.discardRequest,
		Deleted:        suite.formDiscard.Deleted,
	})

	suite.Nil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *discardTestSuite) TestDiscard_Fail() {
	var e error
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.formDiscard.FormSentId, suite.formDiscard.Deleted).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	err := discard(&DiscardInput{
		Db:             suite.gormDb,
		Err:            e,
		DiscardRequest: suite.discardRequest,
		Deleted:        suite.formDiscard.Deleted,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *discardTestSuite) TestDiscard_FailWithErrorInput() {
	var e error = fmt.Errorf("some error")

	err := discard(&DiscardInput{
		Db:             suite.gormDb,
		Err:            e,
		DiscardRequest: suite.discardRequest,
		Deleted:        suite.formDiscard.Deleted,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

type submitTestSuite struct {
	suite.Suite
	sqlFormSubmitInsert      string
	sqlInputSubmissionInsert string
	gormDb                   *gorm.DB
	mockDb                   sqlmock.Sqlmock
	inputSubmissionRequest   request.InputSubmissionRequest
}

func (suite *submitTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *submitTestSuite) SetupSuite() {
	const formSentId = 1
	suite.sqlFormSubmitInsert = `INSERT INTO "form_submissions" ("form_sent_id") VALUES ($1) RETURNING "id"`
	suite.sqlInputSubmissionInsert = `INSERT INTO "input_submissions" ("form_submission_id","input_id","response") VALUES ($1,$2,$3) RETURNING "id"`
	suite.inputSubmissionRequest = request.InputSubmissionRequest{
		FormSentId: formSentId,
		SubmitData: []request.InputData{
			{
				ID:       1,
				Response: "This is a test",
			},
		},
	}
}

func (suite *submitTestSuite) TestSubmit_Success() {
	var e error
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlFormSubmitInsert).WithArgs(suite.inputSubmissionRequest.FormSentId).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInputSubmissionInsert).WithArgs(1, suite.inputSubmissionRequest.SubmitData[0].ID, suite.inputSubmissionRequest.SubmitData[0].Response).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	err := submit(&SubmitInput{
		Db:                     suite.gormDb,
		Err:                    e,
		InputSubmissionRequest: suite.inputSubmissionRequest,
	})

	suite.Nil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *submitTestSuite) TestSubmit_FailOnCreateFormSubmit() {
	var e error
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlFormSubmitInsert).WithArgs(suite.inputSubmissionRequest.FormSentId).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	err := submit(&SubmitInput{
		Db:                     suite.gormDb,
		Err:                    e,
		InputSubmissionRequest: suite.inputSubmissionRequest,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *submitTestSuite) TestSubmit_FailOnCreateInputSubmission() {
	var e error
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlFormSubmitInsert).WithArgs(suite.inputSubmissionRequest.FormSentId).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInputSubmissionInsert).WithArgs(1, suite.inputSubmissionRequest.SubmitData[0].ID, suite.inputSubmissionRequest.SubmitData[0].Response).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	err := submit(&SubmitInput{
		Db:                     suite.gormDb,
		Err:                    e,
		InputSubmissionRequest: suite.inputSubmissionRequest,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *submitTestSuite) TestSubmit_FailWithErrorInput() {
	var e error = fmt.Errorf("some error")
	err := submit(&SubmitInput{
		Db:                     suite.gormDb,
		Err:                    e,
		InputSubmissionRequest: suite.inputSubmissionRequest,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

/* Execute Suites */

func TestDiscardTestSuite(t *testing.T) {
	suite.Run(t, new(discardTestSuite))
}

func TestSubmitTestSuite(t *testing.T) {
	suite.Run(t, new(submitTestSuite))
}
