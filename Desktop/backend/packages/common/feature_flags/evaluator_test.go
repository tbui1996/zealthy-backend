package featureflags

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FlagEvaluatorTestSuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
}

func (suite *FlagEvaluatorTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (s *FlagEvaluatorTestSuite) Test__EvaluatorError() {
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."deleted_at" IS NULL`).
		WillReturnError(fmt.Errorf("error!"))

	evaluator := &FlagEvaluator{
		db: s.gormDb,
	}

	results, err := evaluator.Evaluate()

	s.NotNil(err)
	s.Nil(results)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FlagEvaluatorTestSuite) Test__EvaluatorSuccess() {
	now := time.Now()

	rows := sqlmock.
		NewRows([]string{"id", "name", "key", "is_enabled", "deleted_at", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(2, "name", "key", false, nil, now, "user-id", now, "user-id").
		AddRow(1, "name", "key2", true, nil, now, "user-id", now, "user-id")

	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."deleted_at" IS NULL`).
		WillReturnRows(rows)

	evaluator := &FlagEvaluator{
		db: s.gormDb,
	}

	results, err := evaluator.Evaluate()

	s.Nil(err)
	s.Equal(results.FlagOrDefault("key", true), false)
	s.Equal(results.FlagOrDefault("key2", false), true)
	s.Equal(results.FlagOrDefault("boo2", true), true)
	s.Equal(results.FlagOrDefault("boo", false), false)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFlagEvaluator(t *testing.T) {
	suite.Run(t, new(FlagEvaluatorTestSuite))
}
