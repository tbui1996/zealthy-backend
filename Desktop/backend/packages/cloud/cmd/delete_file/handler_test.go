package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
)

type AddDeleteDateTestSuite struct {
	suite.Suite
}

func (suite *AddDeleteDateTestSuite) TestAddDeleteDate_Success() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	file := model.File{
		ID:               1,
		FileID:           "test id",
		FileName:         "test name",
		FileMimetype:     "image/png",
		SendUserID:       "Okta_00test",
		ChatID:           "chat id",
		DateUploaded:     time.Now(),
		DateLastAccessed: time.Now(),
		FilePath:         "http://s3-test",
		DeletedAt:        nil,
	}

	sqlmock.NewRows([]string{"id", "file_id", "file_name", "file_mimetype", "send_user_id", "chat_id", "date_uploaded", "date_last_accessed", "file_path", "deleted_at"}).AddRow(file.ID, file.FileID, file.FileName, file.FileMimetype, file.SendUserID, file.ChatID, file.DateUploaded, file.DateLastAccessed, file.FilePath, file.DeletedAt)

	date := time.Now()
	fileID := 1
	const sqlUpdate = `UPDATE "files" SET "deleted_at"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, fileID).WillReturnResult(sqlmock.NewResult(0, 1)) // 0 new rows, 1 row updated
	mock.ExpectCommit()

	actual := handler(&AddDeleteDateInput{
		Date:   &date,
		Db:     gdb,
		FileID: fileID,
	})

	suite.Nil(actual)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func (suite *AddDeleteDateTestSuite) TestAddDeleteDate_Fail() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	file := model.File{
		ID:               1,
		FileID:           "test id",
		FileName:         "test name",
		FileMimetype:     "image/png",
		SendUserID:       "Okta_00test",
		ChatID:           "chat id",
		DateUploaded:     time.Now(),
		DateLastAccessed: time.Now(),
		FilePath:         "http://s3-test",
		DeletedAt:        nil,
	}

	sqlmock.NewRows([]string{"id", "file_id", "file_name", "file_mimetype", "send_user_id", "chat_id", "date_uploaded", "date_last_accessed", "file_path", "deleted_at"}).AddRow(file.ID, file.FileID, file.FileName, file.FileMimetype, file.SendUserID, file.ChatID, file.DateUploaded, file.DateLastAccessed, file.FilePath, file.DeletedAt)

	date := time.Now()
	fileID := 1
	const sqlUpdate = `UPDATE "files" SET "deleted_at"=$1 WHERE id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(sqlUpdate).WithArgs(&date, fileID).WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	actualErr := handler(&AddDeleteDateInput{
		Date:   &date,
		Db:     gdb,
		FileID: fileID,
	})

	suite.NotNil(actualErr)
	suite.Error(actualErr)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestAddDeleteDate(t *testing.T) {
	suite.Run(t, new(AddDeleteDateTestSuite))
}
