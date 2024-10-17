package patients

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type PatientRepositorySuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
}

func (suite *PatientRepositorySuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (suite *PatientRepositorySuite) Test_FindSuccess() {
	row := sqlmock.
		NewRows([]string{"id", "name", "last_name", "address", "insurance_id", "birthday", "provider_id"}).
		AddRow(1, "Test", "User", "123 Address", "777777777", time.Now(), "1")

	suite.mock.ExpectQuery(`SELECT * FROM "patients" WHERE "patients"."id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnRows(row)

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	patient, err := repo.Find(model.Patient{ID: 1})

	suite.NoError(err)
	suite.Equal(1, patient.ID)

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *PatientRepositorySuite) Test_FindError() {
	suite.mock.ExpectQuery(`SELECT * FROM "patients" WHERE "patients"."id" = $1 LIMIT 1`).
		WithArgs(1).
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	patient, err := repo.Find(model.Patient{ID: 1})

	suite.Error(err)
	suite.Nil(patient)

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *PatientRepositorySuite) Test_FindAllNoFilterSuccess() {
	row := sqlmock.
		NewRows([]string{"id", "name", "last_name", "address", "insurance_id", "birthday", "provider_id"}).
		AddRow(1, "Test", "User", "123 Address", "777777777", time.Now(), "1").
		AddRow(2, "Cool", "Person", "124 Address", "123456789", time.Now(), "1")

	suite.mock.ExpectQuery(`SELECT * FROM "patients"`).
		WillReturnRows(row)

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	patients, err := repo.FindAll(nil)

	suite.NoError(err)
	suite.Equal(2, len(patients))

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *PatientRepositorySuite) Test_FindAllFilterSuccess() {
	row := sqlmock.
		NewRows([]string{"id", "name", "last_name", "address", "insurance_id", "birthday", "provider_id"}).
		AddRow(1, "Test", "User", "123 Address", "777777777", time.Now(), "1").
		AddRow(2, "Cool", "Person", "124 Address", "123456789", time.Now(), "1")

	suite.mock.ExpectQuery(`SELECT * FROM "patients" WHERE "patients"."provider_id" = $1`).
		WillReturnRows(row)

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	filter := model.Patient{ProviderId: "1"}
	patients, err := repo.FindAll(filter)

	suite.NoError(err)
	suite.Equal(2, len(patients))

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *PatientRepositorySuite) Test_FindAllFailFilter() {
	suite.mock.ExpectQuery(`SELECT * FROM "patients" WHERE "patients"."provider_id" = $1`).
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	filter := model.Patient{ProviderId: "1"}
	patients, err := repo.FindAll(filter)

	suite.Error(err)
	suite.Nil(patients)

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *PatientRepositorySuite) Test_FindAllFailNoFilter() {
	suite.mock.ExpectQuery(`SELECT * FROM "patients"`).
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	repo := &PatientRepository{
		DB: suite.gormDb,
	}

	patients, err := repo.FindAll(nil)

	suite.Error(err)
	suite.Nil(patients)

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPatientRepo(t *testing.T) {
	suite.Run(t, new(PatientRepositorySuite))
}
