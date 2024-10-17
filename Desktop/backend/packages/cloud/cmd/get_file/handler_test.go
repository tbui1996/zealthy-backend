package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GetFileTestSuite struct {
	suite.Suite
}

func (suite *GetFileTestSuite) TestGetFile_Success() {
	var mockFiles []model.File
	mockFile := model.File{
		ID:               1,
		FileID:           "uuid.csv",
		FileName:         "test.csv",
		FileMimetype:     ".csv",
		SendUserID:       "circulohealth",
		ChatID:           "circulator",
		DateUploaded:     time.Now(),
		DateLastAccessed: time.Now(),
		FilePath:         "/test/path/",
		MemberID:         "",
		DeletedAt:        nil,
	}
	mockFiles = append(mockFiles, mockFile)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	mockRow := sqlmock.
		NewRows([]string{"id", "file_id", "file_name", "file_mimetype", "send_user_id", "chat_id", "date_uploaded", "date_last_accessed", "file_path", "member_id", "deleted_at"}).
		AddRow(mockFiles[0].ID, mockFiles[0].FileID, mockFiles[0].FileName, mockFiles[0].FileMimetype, mockFiles[0].SendUserID, mockFiles[0].ChatID, mockFiles[0].DateUploaded, mockFiles[0].DateLastAccessed, mockFiles[0].FilePath, mockFiles[0].MemberID, mockFiles[0].DeletedAt)

	mock.ExpectQuery("SELECT(.*)").WillReturnRows(mockRow)

	actualFiles, err := handler(gdb)

	suite.Equal(mockFiles, actualFiles)
	suite.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func (suite *GetFileTestSuite) TestGetFile_Fail() {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	mock.ExpectQuery("SELECT(.*)").WillReturnError(fmt.Errorf("Some error"))

	actualFiles, err := handler(gdb)

	suite.Nil(actualFiles)
	suite.Error(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetFile(t *testing.T) {
	suite.Run(t, new(GetFileTestSuite))
}
