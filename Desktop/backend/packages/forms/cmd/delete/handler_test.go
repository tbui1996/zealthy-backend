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

type FindFormItemTestSuite struct {
	suite.Suite
}

func (suite *FindFormItemTestSuite) TestFindFormItem_ShouldReturnFormToDeleteIfCallFirst() {
	db, mock, _ := sqlmock.New()

	// Open our mock db connection
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	// Expected result when we run FindFormItem()
	form := model.Form{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		Created:     time.Now(),
		Creator:     "John Smith",
		CreatorId:   "123",
		DeletedAt:   nil,
	}

	// Insert mock row in our mock db
	row := sqlmock.NewRows([]string{"id", "title", "description", "created", "creator_id", "creator", "deleted_at"}).AddRow(form.ID, form.Title, form.Description, form.Created, form.CreatorId, form.Creator, form.DeletedAt)

	formID := "1"
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

	formID := "1"
	mock.ExpectQuery(
		"SELECT(.*)").
		WithArgs(formID).
		WillReturnRows(sqlmock.NewRows(nil))

	actual, actualErr := findFormItem(formID, gdb)

	suite.NotNil(actualErr)
	suite.Empty(actual)
}

type AddDeleteDateTestSuite struct {
	suite.Suite
}

func (suite *AddDeleteDateTestSuite) TestAddDeleteDate_ShouldReturnNilIfCallAddDeleteDate() {
	// QueryMatcherEqual will do a full case sensitive match
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

	date := time.Now()
	formID := "1"
	const sqlUpdate = `UPDATE "forms" SET "deleted_at"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, formID).WillReturnResult(sqlmock.NewResult(0, 1)) // 0 new rows, 1 row updated
	mock.ExpectCommit()

	actual := addDeleteDate(&AddDeleteDateInput{
		Date:   &date,
		Db:     gdb,
		FormID: formID,
	})

	suite.Nil(actual)
}

func (suite *AddDeleteDateTestSuite) TestAddDeleteDate_ShouldReturnErrorIfCallAddDeleteDateErrors() {
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

	date := time.Now()
	formID := "1"
	const sqlUpdate = `UPDATE "forms" SET "deleted_at"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, formID).WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()
	actualErr := addDeleteDate(&AddDeleteDateInput{
		Date:   &date,
		Db:     gdb,
		FormID: formID,
	})

	suite.NotNil(actualErr)
	suite.Error(actualErr)
}

/* Execute Suites */

func TestFindFormItemTestSuite(t *testing.T) {
	suite.Run(t, new(FindFormItemTestSuite))
}

func TestAddDeleteDateTestSuite(t *testing.T) {
	suite.Run(t, new(AddDeleteDateTestSuite))
}
