package main

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestInput struct {
	HandlerInput HandlerInput
	MockDB       sqlmock.Sqlmock
}

func TestHandler(t *testing.T) {
	var username = "test_user"
	var poolID = "pool_id"
	organizationName := "test"
	organizationId := 0

	inputBuilder := func(dbErr error) TestInput {
		db, mock, _ := sqlmock.New()
		gdb, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

		row := sqlmock.
			NewRows([]string{"name"}).
			AddRow(organizationName)

		if dbErr != nil {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "external_user_organizations" WHERE name = $1 ORDER BY "external_user_organizations"."id" LIMIT 1`)).WithArgs(organizationName).WillReturnRows(row)

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "external_users" ("id","external_user_organization_id") VALUES ($1,$2)`)).WithArgs(username, organizationId).WillReturnError(dbErr)
			mock.ExpectRollback()
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "external_user_organizations" WHERE name = $1 ORDER BY "external_user_organizations"."id" LIMIT 1`)).WithArgs(organizationName).WillReturnRows(row)

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "external_users" ("id","external_user_organization_id") VALUES ($1,$2)`)).WithArgs(username, organizationId).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		return TestInput{HandlerInput: HandlerInput{UserName: username, PoolID: poolID, DB: gdb, Logger: zaptest.NewLogger(t), OrganizationName: organizationName}, MockDB: mock}
	}

	tests := []struct {
		input       TestInput
		expectedErr *exception.SonarError
	}{
		{
			input:       inputBuilder(nil),
			expectedErr: nil,
		},
		{
			input:       inputBuilder(errors.New("dbErr")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "dbErr"),
		},
		{
			input:       inputBuilder(nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "idpErr"),
		},
		{
			input:       inputBuilder(errors.New("dbErr")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "idpErr, dbErr"),
		},
	}

	for _, test := range tests {
		actualErr := handler(test.input.HandlerInput)

		if actualErr != nil {
			assert.Equal(t, test.expectedErr.StatusCode, actualErr.StatusCode)
		}

		if err := test.input.MockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
