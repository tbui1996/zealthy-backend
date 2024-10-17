package data

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/aws"
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type FeatureFlagRepositorySuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
	userID string
}

func (suite *FeatureFlagRepositorySuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
	suite.userID = "test-user-id"
}
func (s *FeatureFlagRepositorySuite) Test__SaveNew() {
	row := sqlmock.
		NewRows([]string{"id"}).
		AddRow(1)

	flagModel := model.NewFeatureFlagWithUserId(model.FeatureFlag{
		Key:  "FlagKey",
		Name: "Name",
	}, &s.userID)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "feature_flags"."flags" ("key","name","is_enabled","created_by","updated_by","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(flagModel.Key, flagModel.Name, flagModel.IsEnabled, *flagModel.CreatedBy, *flagModel.UpdatedBy, flagModel.DeletedAt).
		WillReturnRows(row)
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Save(flagModel)

	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__SaveExisting() {
	now := time.Now()

	flagModel := model.NewFeatureFlagWithUserId(model.FeatureFlag{
		Key:  "FlagKey",
		Name: "Name",
	}, &s.userID)

	flagModel.CreatedAt = &now
	flagModel.Id = 1

	sqlUpdate := `UPDATE "feature_flags"."flags" SET "name"=$1,"is_enabled"=$2,"updated_at"=$3,"updated_by"=$4,"deleted_at"=$5 WHERE "id" = $6 AND "flags"."deleted_at" IS NULL`

	s.mock.ExpectBegin()
	s.mock.ExpectExec(sqlUpdate).
		WithArgs(flagModel.Name, flagModel.IsEnabled, AnyTime{}, *flagModel.UpdatedBy, flagModel.DeletedAt, flagModel.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Save(flagModel)

	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__SaveDuplicate() {
	flagModel := model.NewFeatureFlagWithUserId(model.FeatureFlag{
		Key:  "FlagKey",
		Name: "Name",
	}, &s.userID)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "feature_flags"."flags" ("key","name","is_enabled","created_by","updated_by","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(flagModel.Key, flagModel.Name, flagModel.IsEnabled, *flagModel.CreatedBy, *flagModel.UpdatedBy, flagModel.DeletedAt).
		WillReturnError(fmt.Errorf(`"flags_key_key" (SQLSTATE 23505)`))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Save(flagModel)

	s.Equal(flagerror.KEY_CONFLICT, err.Code())

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__SaveError() {
	flagModel := model.NewFeatureFlagWithUserId(model.FeatureFlag{
		Key:  "FlagKey",
		Name: "Name",
	}, &s.userID)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "feature_flags"."flags" ("key","name","is_enabled","created_by","updated_by","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(flagModel.Key, flagModel.Name, flagModel.IsEnabled, *flagModel.CreatedBy, *flagModel.UpdatedBy, flagModel.DeletedAt).
		WillReturnError(fmt.Errorf("unknown"))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Save(flagModel)

	s.NotNil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__FindAllError() {
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."deleted_at" IS NULL`).
		WillReturnError(fmt.Errorf("error!"))

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	results, err := repo.FindAll()

	s.NotNil(err)
	s.Nil(results)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__FindAllSuccess() {
	now := time.Now()
	expected := model.FeatureFlag{
		Id:        1,
		Name:      "name",
		Key:       "key",
		IsEnabled: false,
		CreatedAt: &now,
		CreatedBy: aws.String("user-id"),
		UpdatedAt: &now,
		UpdatedBy: aws.String("user-id"),
	}
	rows := sqlmock.
		NewRows([]string{"id", "name", "key", "is_enabled", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(1, "name", "key", false, now, "user-id", now, "user-id")
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."deleted_at" IS NULL`).
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	results, err := repo.FindAll()

	s.Nil(err)
	s.Equal(expected, (*results)[0])

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__FindSuccess() {
	now := time.Now()
	expected := model.FeatureFlag{
		Id:        2,
		Name:      "name",
		Key:       "key",
		IsEnabled: false,
		CreatedAt: &now,
		CreatedBy: aws.String("user-id"),
		UpdatedAt: &now,
		UpdatedBy: aws.String("user-id"),
	}
	rows := sqlmock.
		NewRows([]string{"id", "name", "key", "is_enabled", "deleted_at", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(2, "name", "key", false, nil, now, "user-id", now, "user-id")
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."id" = $1 AND "flags"."deleted_at" IS NULL ORDER BY "flags"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	result, err := repo.Find(2)

	s.Nil(err)
	s.Equal(expected, *result)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__FindFailNotFound() {
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."id" = $1 AND "flags"."deleted_at" IS NULL ORDER BY "flags"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnError(gorm.ErrRecordNotFound)

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	result, err := repo.Find(2)

	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__FindFail() {
	s.mock.ExpectQuery(`SELECT * FROM "feature_flags"."flags" WHERE "flags"."id" = $1 AND "flags"."deleted_at" IS NULL ORDER BY "flags"."id" LIMIT 1`).
		WithArgs(2).
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	result, err := repo.Find(2)

	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__DeleteSuccess() {
	now := time.Now()
	expected := model.FeatureFlag{
		Id:        2,
		Name:      "name",
		Key:       "key",
		IsEnabled: false,
		CreatedAt: &now,
		CreatedBy: aws.String("user-id"),
		UpdatedAt: &now,
		UpdatedBy: aws.String("user-id"),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "feature_flags"."flags" SET "deleted_at"=$1 WHERE "flags"."id" = $2 AND "flags"."id" = $3 AND "flags"."deleted_at" IS NULL`).
		WithArgs(AnyTime{}, 2, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Delete(&expected)

	s.Nil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__DeleteFail() {
	now := time.Now()
	expected := model.FeatureFlag{
		Id:        2,
		Name:      "name",
		Key:       "key",
		IsEnabled: false,
		CreatedAt: &now,
		CreatedBy: aws.String("user-id"),
		UpdatedAt: &now,
		UpdatedBy: aws.String("user-id"),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "feature_flags"."flags" SET "deleted_at"=$1 WHERE "flags"."id" = $2 AND "flags"."id" = $3 AND "flags"."deleted_at" IS NULL`).
		WithArgs(AnyTime{}, 2, 2).
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Delete(&expected)

	s.NotNil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *FeatureFlagRepositorySuite) Test__DeleteFailNotFound() {
	now := time.Now()
	expected := model.FeatureFlag{
		Id:        2,
		Name:      "name",
		Key:       "key",
		IsEnabled: false,
		CreatedAt: &now,
		CreatedBy: aws.String("user-id"),
		UpdatedAt: &now,
		UpdatedBy: aws.String("user-id"),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "feature_flags"."flags" SET "deleted_at"=$1 WHERE "flags"."id" = $2 AND "flags"."id" = $3 AND "flags"."deleted_at" IS NULL`).
		WithArgs(AnyTime{}, 2, 2).
		WillReturnError(gorm.ErrRecordNotFound)
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &FeatureFlagRepository{
		DB:     s.gormDb,
		UserID: s.userID,
		Logger: logger,
	}

	err := repo.Delete(&expected)

	s.NotNil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFeatureFlagRepo(t *testing.T) {
	suite.Run(t, new(FeatureFlagRepositorySuite))
}
