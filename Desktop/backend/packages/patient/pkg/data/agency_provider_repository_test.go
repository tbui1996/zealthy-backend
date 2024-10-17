package data

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AgencyProviderRepositorySuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
}

func (suite *AgencyProviderRepositorySuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (s *AgencyProviderRepositorySuite) Test__FindAllSuccess() {
	now := time.Now()
	expected := model.AgencyProvider{
		AgencyProviderId:      "12",
		DoddNumber:            "13",
		NationalProviderId:    "1",
		FirstName:             "a",
		MiddleName:            "b",
		LastName:              "c",
		Suffix:                "d",
		BusinessName:          "here",
		BusinessTIN:           "3",
		BusinessAddress1:      "circulo",
		BusinessAddress2:      "address2",
		BusinessCity:          "columbus",
		BusinessState:         "OH",
		BusinessZip:           "00000",
		CreatedTimestamp:      &now,
		LastModifiedTimestamp: &now,
	}
	rows := sqlmock.
		NewRows([]string{"agency_provider_id", "dodd_number", "national_provider_id", "first_name",
			"middle_name", "last_name", "suffix",
			"business_name", "business_tin", "business_address_1",
			"business_address_2", "business_city", "business_state",
			"business_zip", "created_timestamp", "last_modified_timestamp"}).
		AddRow("12",
			"13", "1", "a",
			"b", "c", "d",
			"here", "3", "circulo",
			"address2", "columbus", "OH",
			"00000", now, now)
	s.mock.ExpectQuery(`SELECT * FROM "agency_provider"`).
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	results, err := repo.FindAll()

	s.Nil(err)
	s.Equal(expected, (*results)[0])

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__FindAllError() {
	s.mock.ExpectQuery(`SELECT * FROM "agency_provider"`).WillReturnError(fmt.Errorf("error!"))

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	results, err := repo.FindAll()

	s.NotNil(err)
	s.Nil(results)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

}

func (s *AgencyProviderRepositorySuite) Test__FindSuccess() {
	expected := model.AgencyProvider{
		AgencyProviderId:   "123789",
		DoddNumber:         "2",
		NationalProviderId: "3",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessCity:       "columbus",
		BusinessState:      "oh",
		BusinessZip:        "00000",
	}
	rows := sqlmock.
		NewRows([]string{"agency_provider_id", "dodd_number", "national_provider_id",
			"first_name", "middle_name", "last_name",
			"suffix", "business_name", "business_tin",
			"business_address_1", "business_address_2", "business_city",
			"business_state", "business_zip"}).
		AddRow("123789", "2", "3", "thomas", "", "bui", "", "circulo", "", "here", "", "columbus", "oh", "00000")
	s.mock.ExpectQuery(`SELECT * FROM "agency_provider" WHERE agency_provider_id = $1 ORDER BY "agency_provider"."agency_provider_id" LIMIT 1`).
		WithArgs("123789").
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("123789")

	s.Nil(err)
	s.Equal(expected, *result)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__FindFailNotFound() {
	s.mock.ExpectQuery(`SELECT * FROM "agency_provider" WHERE agency_provider_id = $1 ORDER BY "agency_provider"."agency_provider_id" LIMIT 1`).
		WithArgs("2").
		WillReturnError(gorm.ErrRecordNotFound)

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("2")

	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__FindFail() {
	s.mock.ExpectQuery(`SELECT * FROM "agency_provider" WHERE agency_provider_id = $1 ORDER BY "agency_provider"."agency_provider_id" LIMIT 1`).
		WithArgs("2").WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("2")

	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

}

func (s *AgencyProviderRepositorySuite) Test__SaveNew() {

	agencyProvider := model.NewAgencyProvider(model.AgencyProvider{
		DoddNumber:         "2",
		NationalProviderId: "3",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "ohio",
		BusinessAddress2:   "",
		BusinessCity:       "ohio",
		BusinessState:      "ohio",
		BusinessZip:        "ohio",
	})

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT INTO "agency_provider" ("dodd_number","national_provider_id","first_name","middle_name","last_name","suffix","business_name","business_tin","business_address_1","business_address_2","business_city","business_state","business_zip") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`).
		WithArgs(agencyProvider.DoddNumber, agencyProvider.NationalProviderId, agencyProvider.FirstName, agencyProvider.MiddleName, agencyProvider.LastName, agencyProvider.Suffix, agencyProvider.BusinessName, agencyProvider.BusinessTIN, agencyProvider.BusinessAddress1, agencyProvider.BusinessAddress2, agencyProvider.BusinessCity, agencyProvider.BusinessState, agencyProvider.BusinessZip).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Save(agencyProvider)

	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__SaveExisting() {
	now := time.Now()

	agencyProvider := model.NewAgencyProvider(model.AgencyProvider{
		DoddNumber:         "2",
		NationalProviderId: "3",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "ohio",
		BusinessAddress2:   "",
		BusinessCity:       "ohio",
		BusinessState:      "ohio",
		BusinessZip:        "ohio",
	})
	agencyProvider.AgencyProviderId = "12345"
	agencyProvider.DoddNumber = "3"
	agencyProvider.CreatedTimestamp = &now
	agencyProvider.LastModifiedTimestamp = &now

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "agency_provider" SET "dodd_number"=$1,"national_provider_id"=$2,"first_name"=$3,"middle_name"=$4,"last_name"=$5,"suffix"=$6,"business_name"=$7,"business_tin"=$8,"business_address_1"=$9,"business_address_2"=$10,"business_city"=$11,"business_state"=$12,"business_zip"=$13,"last_modified_timestamp"=$14 WHERE "agency_provider_id" = $15`).
		WithArgs(agencyProvider.DoddNumber, agencyProvider.NationalProviderId, agencyProvider.FirstName, agencyProvider.MiddleName, agencyProvider.LastName, agencyProvider.Suffix, agencyProvider.BusinessName, agencyProvider.BusinessTIN, agencyProvider.BusinessAddress1, agencyProvider.BusinessAddress2, agencyProvider.BusinessCity, agencyProvider.BusinessState, agencyProvider.BusinessZip, agencyProvider.LastModifiedTimestamp, agencyProvider.AgencyProviderId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Save(agencyProvider)

	s.Nil(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__SaveError() {
	agencyProvider := model.NewAgencyProvider(model.AgencyProvider{
		DoddNumber:         "2",
		NationalProviderId: "3",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "ohio",
		BusinessAddress2:   "",
		BusinessCity:       "ohio",
		BusinessState:      "ohio",
		BusinessZip:        "ohio",
	})

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT INTO "agency_provider" ("dodd_number","national_provider_id","first_name","middle_name","last_name","suffix","business_name","business_tin","business_address_1","business_address_2","business_city","business_state","business_zip") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`).
		WithArgs(agencyProvider.DoddNumber, agencyProvider.NationalProviderId, agencyProvider.FirstName, agencyProvider.MiddleName, agencyProvider.LastName, agencyProvider.Suffix, agencyProvider.BusinessName, agencyProvider.BusinessTIN, agencyProvider.BusinessAddress1, agencyProvider.BusinessAddress2, agencyProvider.BusinessCity, agencyProvider.BusinessState, agencyProvider.BusinessZip).
		WillReturnError(fmt.Errorf("unknown"))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &AgencyProviderRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Save(agencyProvider)

	s.Equal(appointmenterror.UNKNOWN, err.Code())

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAgencyProviderRepo(t *testing.T) {
	suite.Run(t, new(AgencyProviderRepositorySuite))
}
