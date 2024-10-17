package main

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FindFormItemTestSuite struct {
	suite.Suite
}

func (suite *FindFormItemTestSuite) TestFindFormItem_ShouldReturnFormToDeleteIfCallFirst() {
	db, mock, _ := sqlmock.New()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	form := model.Form{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		Created:     time.Now(),
		Creator:     "John Smith",
		CreatorId:   "123",
		DeletedAt:   nil,
	}

	row := sqlmock.NewRows([]string{"id", "title", "description", "created", "creator_id", "creator", "deleted_at"}).AddRow(form.ID, form.Title, form.Description, form.Created, form.CreatorId, form.Creator, form.DeletedAt)

	formID := 1
	mock.ExpectQuery(
		"SELECT(.*)").
		WithArgs(formID).
		WillReturnRows(row)

	actual, actualErr := findFormItem(formID, gdb)

	suite.Nil(actualErr)
	suite.Equal(form, actual)
}

func (suite *FindFormItemTestSuite) TestFindFormItem_ShouldReturnErrorIfCallFirstError() {
	db, mock, _ := sqlmock.New()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	formID := 1
	mock.ExpectQuery(
		"SELECT(.*)").
		WithArgs(formID).
		WillReturnRows(sqlmock.NewRows(nil))

	actual, actualErr := findFormItem(formID, gdb)

	suite.NotNil(actualErr)
	suite.Empty(actual)
}

type EditFormItemTestSuite struct {
	suite.Suite
}

func (suite *EditFormItemTestSuite) TestEditFormItem_ShouldReturnNilIfCallEditFormItem() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	form := model.Form{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		Created:     time.Now(),
		Creator:     "John Smith",
		CreatorId:   "123",
		DeletedAt:   nil,
	}

	sqlmock.NewRows([]string{"id", "title", "description", "created", "creator_id", "creator", "deleted_at"}).AddRow(form.ID, form.Title, form.Description, form.Created, form.CreatorId, form.Creator, form.DeletedAt)

	request := request.EditForm{
		ID:          1,
		Title:       "Update test",
		Description: "Update Description",
	}

	const sqlUpdate = `UPDATE "forms" SET "title"=$1,"description"=$2,"created"=$3,"creator"=$4,"creator_id"=$5,"deleted_at"=$6,"date_closed"=$7 WHERE "id" = $8`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(request.Title, request.Description, form.Created, form.Creator, form.CreatorId, form.DeletedAt, nil, form.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	actual := editFormItem(&EditFormItemInput{
		Form: form,
		Req:  request,
		Db:   gdb,
	})

	suite.Nil(actual)
}

/* Execute Suites */

func TestFindFormItemTestSuite(t *testing.T) {
	suite.Run(t, new(FindFormItemTestSuite))
}
func TestEditFormItemTestSuite(t *testing.T) {
	suite.Run(t, new(EditFormItemTestSuite))
}
