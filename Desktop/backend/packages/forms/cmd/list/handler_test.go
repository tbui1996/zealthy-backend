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

type GetAllFormsTestSuite struct {
	suite.Suite
	gormDb   *gorm.DB
	mockDb   sqlmock.Sqlmock
	sqlQuery string
}

func (suite *GetAllFormsTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *GetAllFormsTestSuite) SetupSuite() {
	suite.sqlQuery = `SELECT * FROM form.forms forms LEFT JOIN (SELECT ROW_NUMBER() OVER (PARTITION BY form.form_sents.form_id ORDER BY form.form_sents.sent desc) AS rownum,form.form_sents.sent, form.form_sents.form_id FROM form.form_sents) sents ON sents.form_id = forms.id AND sents.rownum = 1 WHERE forms.deleted_at IS NULL`

}

func (suite *GetAllFormsTestSuite) TestGetAllForms_Success() {
	form := model.Form{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		Created:     time.Now(),
		Creator:     "John Smith",
		CreatorId:   "123",
		DeletedAt:   nil,
	}

	var formResult = []Result{
		{
			ID:          form.ID,
			Title:       form.Title,
			Description: form.Description,
			Created:     form.Created,
			Creator:     form.Creator,
			Sent:        time.Time{},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "created", "creator_id", "creator", "deleted_at"}).AddRow(form.ID, form.Title, form.Description, form.Created, form.CreatorId, form.Creator, form.DeletedAt)

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnRows(rows)

	actual, actualErr := getAllForms(suite.gormDb)

	suite.Nil(actualErr)
	suite.NotNil(actual)
	suite.EqualValues(formResult, actual)

	err := suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *GetAllFormsTestSuite) TestGetAllForms_Fail() {

	suite.mockDb.ExpectQuery(suite.sqlQuery).WillReturnError(fmt.Errorf("some error"))

	actual, actualErr := getAllForms(suite.gormDb)

	suite.NotNil(actualErr)
	suite.Nil(actual)
	suite.Error(actualErr)

	err := suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

/* Execute Suites */

func TestGetAllFormsTestSuite(t *testing.T) {
	suite.Run(t, new(GetAllFormsTestSuite))
}
