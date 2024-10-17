package main

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FileUploadTestSuite struct {
	suite.Suite
	S3Mock  *mocks.S3API
	Request FileUploadHandler
	SQLMock sqlmock.Sqlmock
	Insert  string
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (suite *FileUploadTestSuite) SetupTest() {
	svc := new(mocks.S3API)

	db, theMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.SQLMock = theMock

	suite.S3Mock = svc

	suite.Insert = `INSERT INTO "files" ("file_id","file_name","file_mimetype","send_user_id","chat_id","date_uploaded","date_last_accessed","file_path","member_id","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`

	suite.Request = FileUploadHandler{
		S3:         svc,
		DB:         gdb,
		Logger:     zaptest.NewLogger(suite.T()),
		BucketName: "test",
		Username:   "test",
		UploadRequest: request.FileUploadRequest{
			ChatId:   "1",
			FileId:   "1",
			Filename: "test.jpg",
		},
	}
}

func (suite *FileUploadTestSuite) TestFileUpload_Success() {
	row := sqlmock.NewRows([]string{"id"}).AddRow(1)

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectQuery(suite.Insert).
		WithArgs("1", "test.jpg", "image/jpeg", "test", "1", AnyTime{}, AnyTime{}, "https://s3.us-east-2.amazonaws.com/test/1", "", nil).
		WillReturnRows(row)
	suite.SQLMock.ExpectCommit()

	res, err := handler(&suite.Request)

	suite.Nil(err)
	suite.NotNil(res)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *FileUploadTestSuite) TestFileUpload_FailOnFileExt() {
	req := suite.Request

	req.UploadRequest.Filename = "test.sdfhgsdfhsfgh"

	res, err := handler(&req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *FileUploadTestSuite) TestFileUpload_FailOnDB() {
	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectQuery(suite.Insert).
		WithArgs("1", "test.jpg", "image/jpeg", "test", "1", AnyTime{}, AnyTime{}, "https://s3.us-east-2.amazonaws.com/test/1", "", nil).
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.SQLMock.ExpectRollback()

	suite.S3Mock.On("DeleteObject", mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)

	res, err := handler(&suite.Request)

	suite.NotNil(err)
	suite.Nil(res)
	suite.S3Mock.AssertCalled(suite.T(), "DeleteObject", mock.Anything)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *FileUploadTestSuite) TestFileUpload_FailOnDelete() {
	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectQuery(suite.Insert).
		WithArgs("1", "test.jpg", "image/jpeg", "test", "1", AnyTime{}, AnyTime{}, "https://s3.us-east-2.amazonaws.com/test/1", "", nil).
		WillReturnError(errors.New("FAKE TEST ERROR, IGNORE"))
	suite.SQLMock.ExpectRollback()

	suite.S3Mock.On("DeleteObject", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	res, err := handler(&suite.Request)

	suite.NotNil(err)
	suite.Nil(res)
	suite.S3Mock.AssertCalled(suite.T(), "DeleteObject", mock.Anything)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *FileUploadTestSuite) TestFileUpload_FailOnMarshal() {
	row := sqlmock.NewRows([]string{"id"})

	suite.SQLMock.ExpectBegin()
	suite.SQLMock.ExpectQuery(suite.Insert).
		WithArgs("1", "test.jpg", "image/jpeg", "test", "1", AnyTime{}, AnyTime{}, "https://s3.us-east-2.amazonaws.com/test/1", "", nil).
		WillReturnRows(row)
	suite.SQLMock.ExpectCommit()

	res, err := handler(&suite.Request)

	suite.NotNil(err)
	suite.Nil(res)

	if err := suite.SQLMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFileUploadTestSuite(t *testing.T) {
	suite.Run(t, new(FileUploadTestSuite))
}
