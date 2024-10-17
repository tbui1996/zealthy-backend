package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SendThroughRouterTestSuite struct {
	suite.Suite
}

func (suite *CreateSendRecordTestSuite) TestSendThroughRouter_Success() {
	client := new(router.MockRouter)

	mockClient := &router.Session{
		Router: client,
	}

	input := &router.RouterSendInput{
		Source:     "forms",
		Action:     "forms",
		Procedure:  "send",
		Body:       "test",
		Recipients: []string{},
	}

	client.On("Send", input).Return(nil)

	err := sendThroughRouter(mockClient, "test")
	suite.Nil(err)
}

func (suite *CreateSendRecordTestSuite) TestSendThroughRouter_Fail() {
	client := new(router.MockRouter)

	mockClient := &router.Session{
		Router: client,
	}

	input := &router.RouterSendInput{
		Source:     "forms",
		Action:     "forms",
		Procedure:  "send",
		Body:       "test",
		Recipients: []string{},
	}

	client.On("Send", input).Return(errors.New("test"))

	err := sendThroughRouter(mockClient, "test")
	suite.NotNil(err)
}

type CreateSendRecordTestSuite struct {
	suite.Suite
	sqlInsert string
	gormDb    *gorm.DB
	mockDb    sqlmock.Sqlmock
	sent      model.FormSent
}

func (suite *CreateSendRecordTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.gormDb = gdb
	suite.mockDb = mock
}

func (suite *CreateSendRecordTestSuite) SetupSuite() {
	suite.sqlInsert = `INSERT INTO "form_sents" ("form_id","sent") VALUES ($1,$2) RETURNING "id"`
	suite.sent = model.FormSent{
		ID:     1,
		FormId: 1,
		Sent:   time.Now(),
	}
}

func (suite *CreateSendRecordTestSuite) AfterTest(suiteName, testName string) {
	err := suite.mockDb.ExpectationsWereMet()
	suite.Nil(err)
}

func (suite *CreateSendRecordTestSuite) TestCreateFormItem_Success() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.sent.FormId, suite.sent.Sent).WillReturnRows(suite.mockDb.NewRows([]string{"id"}).AddRow(1))
	suite.mockDb.ExpectCommit()

	actual, err := createSendRecord(&CreateSendRecordInput{
		FormID: 1,
		Sent:   suite.sent.Sent,
		Db:     suite.gormDb,
	})

	suite.Nil(err)
	suite.EqualValues(actual, suite.sent)
}

func (suite *CreateSendRecordTestSuite) TestCreateFormItem_Fail() {
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(suite.sqlInsert).WithArgs(suite.sent.FormId, suite.sent.Sent).WillReturnError(fmt.Errorf("some error"))
	suite.mockDb.ExpectRollback()

	_, err := createSendRecord(&CreateSendRecordInput{
		FormID: 1,
		Sent:   suite.sent.Sent,
		Db:     suite.gormDb,
	})

	suite.NotNil(err)
}

/* Execute Suites */

func TestSendThroughRouterTestSuite(t *testing.T) {
	suite.Run(t, new(SendThroughRouterTestSuite))
}

func TestCreateSendRecordTestSuite(t *testing.T) {
	suite.Run(t, new(CreateSendRecordTestSuite))
}
