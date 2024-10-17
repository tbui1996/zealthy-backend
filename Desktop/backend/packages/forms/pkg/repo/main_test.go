package repo

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

type RepoInputsTestSuite struct {
	suite.Suite
	sqlQuery   string
	repository *Repository
	mockDb     sqlmock.Sqlmock
	inputs     model.Input
}

func (suite *RepoInputsTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.mockDb = mock
	suite.repository = NewRepository(gdb)
}

func (suite *RepoInputsTestSuite) SetupSuite() {
	suite.inputs = model.Input{
		ID:      1,
		Order:   0,
		Type:    "text",
		FormId:  1,
		Label:   "Enter a text",
		Options: []string{},
	}
	suite.sqlQuery = `SELECT * FROM "inputs" WHERE form_id = $1`
}

func (suite *RepoInputsTestSuite) AfterTest(suiteName, testName string) {
	err := suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *RepoInputsTestSuite) TestInput_Success() {
	const formId = "1"
	rows := sqlmock.NewRows([]string{"id", "order", "type", "label", "options", "form_id"}).AddRow(suite.inputs.ID, suite.inputs.Order, suite.inputs.Type, suite.inputs.Label, suite.inputs.Options, suite.inputs.FormId)

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnRows(rows)

	actual, err := suite.repository.Inputs(formId)

	suite.Nil(err)
	suite.NotNil(actual)
	suite.EqualValues(actual, []model.Input{suite.inputs})
}

func (suite *RepoInputsTestSuite) TestInput_Fail() {
	const formId = "1"
	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnError(fmt.Errorf("some error"))

	_, err := suite.repository.Inputs(formId)

	suite.NotNil(err)
}

type RepoFormTestSuite struct {
	suite.Suite
	sqlQuery   string
	repository *Repository
	mockDb     sqlmock.Sqlmock
	form       model.Form
}

func (suite *RepoFormTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.repository = NewRepository(gdb)
	suite.mockDb = mock
}

func (suite *RepoFormTestSuite) SetupSuite() {
	form := model.Form{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		Created:     time.Now(),
		Creator:     "John Smith",
		CreatorId:   "123",
		DeletedAt:   nil,
	}

	suite.form = form
	suite.sqlQuery = `SELECT * FROM "forms" WHERE "forms"."id" = $1 ORDER BY "forms"."id" LIMIT 1`
}

func (suite *RepoFormTestSuite) AfterTest(suiteName, testName string) {
	err := suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *RepoFormTestSuite) TestForm_Success() {
	const formId = "1"
	rows := sqlmock.NewRows([]string{"id", "title", "description", "created", "creator_id", "creator", "deleted_at"}).AddRow(suite.form.ID, suite.form.Title, suite.form.Description, suite.form.Created, suite.form.CreatorId, suite.form.Creator, suite.form.DeletedAt)

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnRows(rows)

	actual, err := suite.repository.Form(formId)

	suite.Nil(err)
	suite.NotNil(actual)
}

func (suite *RepoFormTestSuite) TestForm_Fail() {
	const formId = "1"
	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnError(fmt.Errorf("some error"))

	_, err := suite.repository.Form(formId)

	suite.NotNil(err)
}

func TestRepoInputsTestSuite(t *testing.T) {
	suite.Run(t, new(RepoInputsTestSuite))
}

func TestRepoFormTestSuite(t *testing.T) {
	suite.Run(t, new(RepoFormTestSuite))
}
