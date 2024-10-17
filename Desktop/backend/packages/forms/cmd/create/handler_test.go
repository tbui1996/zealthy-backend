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

type CreateFormItemTestSuite struct {
	suite.Suite
	newId        int
	sqlInsert    string
	gormDb       *gorm.DB
	mockDb       sqlmock.Sqlmock
	created      time.Time
	form         request.CreateForm
	expectedForm model.Form
}

func (suite *CreateFormItemTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *CreateFormItemTestSuite) SetupSuite() {
	var created = time.Now()
	const newId = 1

	form := request.CreateForm{
		Title:       "test title",
		Description: "test description",
		Creator:     "John Smith",
		CreatorId:   "123",
	}

	expectedForm := model.Form{
		ID:          newId,
		Title:       form.Title,
		Description: form.Description,
		Created:     created,
		Creator:     form.Creator,
		CreatorId:   form.CreatorId,
		DeletedAt:   nil,
		DateClosed:  nil,
	}

	suite.created = created
	suite.newId = newId
	suite.form = form
	suite.expectedForm = expectedForm

	suite.sqlInsert = `INSERT INTO "forms" ("title","description","created","creator","creator_id","deleted_at","date_closed") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
}

func (suite *CreateFormItemTestSuite) TestCreateFormItem_Success() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.form.Title, suite.form.Description, suite.created, suite.form.Creator, suite.form.CreatorId, nil, nil).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(suite.newId))
	suite.mockDb.ExpectCommit()

	actual, err := createFormItem(&CreateFormItemInput{
		In:      suite.form,
		Db:      suite.gormDb,
		Created: suite.created,
	})

	suite.Nil(err)
	suite.NotNil(actual)
	suite.Equal(suite.expectedForm, actual)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *CreateFormItemTestSuite) TestCreateFormItem_Fail() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.form.Title, suite.form.Description, suite.created, suite.form.Creator, suite.form.CreatorId, nil, nil).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	_, err := createFormItem(&CreateFormItemInput{
		In:      suite.form,
		Db:      suite.gormDb,
		Created: suite.created,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

type CreateInputItemsTestSuite struct {
	suite.Suite
	sqlInsert string
	gormDb    *gorm.DB
	mockDb    sqlmock.Sqlmock
	inputs    model.Input
	form      request.CreateForm
}

func (suite *CreateInputItemsTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *CreateInputItemsTestSuite) SetupSuite() {
	const formId = 1

	inputs := model.Input{
		ID:      1,
		Order:   0,
		Type:    "text",
		FormId:  formId,
		Label:   "Enter a text",
		Options: []string{},
	}

	form := request.CreateForm{
		Title:       "test title",
		Description: "test description",
		Creator:     "John Smith",
		CreatorId:   "123",
		Inputs:      []model.Input{inputs},
	}

	suite.inputs = inputs
	suite.form = form
	suite.sqlInsert = `INSERT INTO "inputs" ("order","type","form_id","label","options") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
}

func (suite *CreateInputItemsTestSuite) TestCreateInputItems_Success() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.inputs.Order, suite.inputs.Type, suite.inputs.FormId, suite.inputs.Label, suite.inputs.Options).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	err := createInputItems(&CreateInputItemsInput{
		FormId: 1,
		Db:     suite.gormDb,
		Form:   suite.form,
	})

	suite.Nil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *CreateInputItemsTestSuite) TestCreateInputItems_Fail() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.inputs.Order, suite.inputs.Type, suite.inputs.FormId, suite.inputs.Label, suite.inputs.Options).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	err := createInputItems(&CreateInputItemsInput{
		FormId: 1,
		Db:     suite.gormDb,
		Form:   suite.form,
	})

	suite.NotNil(err)

	err = suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

/* Execute Suites */

func TestCreateFormItemTestSuite(t *testing.T) {
	suite.Run(t, new(CreateFormItemTestSuite))
}

func TestCreateInputItemsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateInputItemsTestSuite))
}
