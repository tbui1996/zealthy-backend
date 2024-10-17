package main

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AssociateFileTestSuite struct {
	suite.Suite
	MockFile model.File
}

type ParseRequestTestSuite struct {
	suite.Suite
}

// This anytime struct and function are necessary since the time values will be off by milliseconds between the mock and the associateFile call.
// [time.Time - 2021-10-21 15:42:58.233702 -0700 PDT m=+0.004017203] does not match actual [time.Time - 2021-10-21 15:42:58.233715 -0700 PDT m=+0.004029778]
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (suite *AssociateFileTestSuite) SetupTest() {
	suite.MockFile = model.File{
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
	}
}

func (suite *AssociateFileTestSuite) TestAssociateFile_Success() {
	fileToAssociate := model.File{
		FileID:   "uuid.csv",
		MemberID: "1",
		FilePath: "test",
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	db2, mock2, _ := sqlmock.New()
	defer db2.Close()

	gdb2, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db2,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "date_last_accessed"=$1,"member_id"=$2 WHERE file_id = $3`)).
		WithArgs(AnyTime{}, fileToAssociate.MemberID, fileToAssociate.FileID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta(`UPDATE "patient" SET "circulo_consent_form_link"=$1,"last_modified_timestamp"=$2,"signed_circulo_consent_form"=$3 WHERE medicaid_id = $4`)).
		WithArgs("test", AnyTime{}, true, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock2.ExpectCommit()

	err := associateFile(&AssociateFileHandler{
		DopplerDb: gdb2,
		SonarDb:   gdb,
		File:      &fileToAssociate,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
	if err := mock2.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func (suite *AssociateFileTestSuite) TestAssociateFile_Fail() {
	fileToAssociate := model.File{
		FileID:   "invalidFileID.csv",
		MemberID: "nowAssociated",
		FilePath: "test",
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	db2, mock2, _ := sqlmock.New()
	defer db2.Close()

	gdb2, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db2,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "date_last_accessed"=$1,"member_id"=$2 WHERE file_id = $3`)).
		WithArgs(AnyTime{}, fileToAssociate.MemberID, fileToAssociate.FileID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	mock2.ExpectBegin()
	mock2.ExpectExec(regexp.QuoteMeta(`UPDATE "patient" SET "circulo_consent_form_link"=$1,"last_modified_timestamp"=$2,"signed_circulo_consent_form"=$3 WHERE medicaid_id = $4`)).
		WithArgs("test", AnyTime{}, true, "nowAssociated").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock2.ExpectRollback()

	err := associateFile(&AssociateFileHandler{
		DopplerDb: gdb2,
		SonarDb:   gdb,
		File:      &fileToAssociate,
		Logger:    zaptest.NewLogger(suite.T()),
	})

	suite.Error(err)

	if err := mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
	if err := mock2.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("There were unfulfilled expectations: %s", err)
	}
}

func (suite *ParseRequestTestSuite) TestParseRequest_Success() {
	mockEvent := events.APIGatewayV2HTTPRequest{
		Body: `{"fileID":"test","memberID":"test"}`,
	}
	expectedFile := &model.File{
		FileID:   "test",
		MemberID: "test",
	}
	actualFile, err := parseRequest(mockEvent)

	suite.Equal(expectedFile, actualFile)
	suite.Nil(err)
}

func (suite *ParseRequestTestSuite) TestParseRequest_Fail() {
	mockEvent := events.APIGatewayV2HTTPRequest{
		Body: `"invalid json"`,
	}

	actualFile, err := parseRequest(mockEvent)

	suite.Nil(actualFile)
	suite.Error(err)
}

func TestAssociateFile(t *testing.T) {
	suite.Run(t, new(AssociateFileTestSuite))
}

func TestParseRequest(t *testing.T) {
	suite.Run(t, new(ParseRequestTestSuite))
}
