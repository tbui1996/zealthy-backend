package data

import (
	"database/sql"
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
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (s *PatientRepositorySuite) Test__FindAllSuccess() {
	now := time.Now()
	expected := model.Patient{
		PatientId:                  "1",
		InsuranceId:                "1",
		FirstName:                  "a",
		MiddleName:                 "b",
		LastName:                   "c",
		Suffix:                     "d",
		DateOfBirth:                now,
		PrimaryLanguage:            "eng",
		PreferredGender:            "M",
		EmailAddress:               "circulo",
		HomePhone:                  "3011000000",
		HomeLivingArrangement:      "yes",
		HomeAddress1:               "crest",
		HomeAddress2:               "view",
		HomeCity:                   "md",
		HomeCounty:                 "md",
		HomeState:                  "md",
		HomeZip:                    "md",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
		CreatedTimestamp:           &now,
		LastModifiedTimestamp:      &now,
	}
	rows := sqlmock.
		NewRows([]string{"patient_id", "insurance_id", "first_name",
			"middle_name", "last_name", "suffix",
			"date_of_birth", "primary_language", "preferred_gender",
			"email_address", "home_phone", "home_living_arrangement",
			"home_address_1", "home_address_2", "home_city",
			"home_county", "home_state", "home_zip",
			"signed_circulo_consent_form", "circulo_consent_form_link", "signed_stationmd_consent_form_link",
			"stationmd_consent_form_link", "completed_go_sheet", "marked_as_active", "created_timestamp",
			"last_modified_timestamp"}).
		AddRow("1", "1", "a", "b", "c", "d", now, "eng", "M", "circulo", "3011000000", "yes", "crest", "view", "md", "md", "md", "md", false, "", false, "", false, false, now, now)
	s.mock.ExpectQuery(`SELECT * FROM "patient"`).
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &PatientRepository{
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

func (s *PatientRepositorySuite) Test__SaveDuplicate() {
	now := time.Now()

	patient := model.NewPatient(model.Patient{
		InsuranceId:                "2",
		FirstName:                  "he",
		MiddleName:                 "oh",
		LastName:                   "won",
		Suffix:                     "Mr",
		DateOfBirth:                now,
		PrimaryLanguage:            "en",
		PreferredGender:            "M",
		EmailAddress:               "hewon",
		HomePhone:                  "301655",
		HomeLivingArrangement:      "townhouse",
		HomeAddress1:               "here",
		HomeAddress2:               "nothere",
		HomeCity:                   "md",
		HomeCounty:                 "moco",
		HomeState:                  "md",
		HomeZip:                    "00001",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "empty",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "empty",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
	})

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT INTO "patient" ("insurance_id","first_name","middle_name","last_name","suffix","date_of_birth","primary_language","preferred_gender","email_address","home_phone","home_living_arrangement","home_address_1","home_address_2","home_city","home_county","home_state","home_zip","signed_circulo_consent_form","circulo_consent_form_link","signed_stationmd_consent_form","stationmd_consent_form_link","completed_go_sheet","marked_as_active") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`).
		WithArgs(patient.InsuranceId, patient.FirstName, patient.MiddleName, patient.LastName,
			patient.Suffix, patient.DateOfBirth, patient.PrimaryLanguage, patient.PreferredGender,
			patient.EmailAddress, patient.HomePhone, patient.HomeLivingArrangement, patient.HomeAddress1,
			patient.HomeAddress2, patient.HomeCity, patient.HomeCounty, patient.HomeState,
			patient.HomeZip, patient.SignedCirculoConsentForm, patient.CirculoConsentFormLink, patient.SignedStationMDConsentForm, patient.StationMDConsentFormLink, patient.CompletedGoSheet,
			patient.MarkedAsActive).
		WillReturnError(fmt.Errorf("(SQLSTATE 23505) - patient with this insurance_id already exists"))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &PatientRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Save(patient)

	s.Equal(appointmenterror.INSURANCE_ID_CONFLICT, err.Code())

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *PatientRepositorySuite) Test__FindSuccess() {
	now := time.Now()

	expected := model.Patient{
		PatientId:                  "1",
		InsuranceId:                "2",
		FirstName:                  "thomas",
		MiddleName:                 "",
		LastName:                   "bui",
		Suffix:                     "",
		DateOfBirth:                now,
		PrimaryLanguage:            "",
		PreferredGender:            "",
		EmailAddress:               "circulo",
		HomePhone:                  "3016558003",
		HomeLivingArrangement:      "",
		HomeAddress1:               "",
		HomeAddress2:               "",
		HomeCity:                   "dog",
		HomeCounty:                 "cat",
		HomeState:                  "mouse",
		HomeZip:                    "67890",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
		CreatedTimestamp:           &now,
		LastModifiedTimestamp:      &now,
	}

	rows := sqlmock.NewRows([]string{"patient_id", "insurance_id", "first_name",
		"middle_name", "last_name", "suffix",
		"date_of_birth", "primary_language", "preferred_gender",
		"email_address", "home_phone", "home_living_arrangement",
		"home_address_1", "home_address_2", "home_city",
		"home_county", "home_state", "home_zip",
		"signed_circulo_consent_form", "circulo_consent_form_link", "signed_stationmd_consent_form_link",
		"stationmd_consent_form_link", "completed_go_sheet", "marked_as_active", "created_timestamp",
		"last_modified_timestamp"}).
		AddRow(expected.PatientId, expected.InsuranceId, expected.FirstName, expected.MiddleName, expected.LastName,
			expected.Suffix, expected.DateOfBirth, expected.PrimaryLanguage, expected.PreferredGender, expected.EmailAddress,
			expected.HomePhone, expected.HomeLivingArrangement, expected.HomeAddress1, expected.HomeAddress2, expected.HomeCity,
			expected.HomeCounty, expected.HomeState, expected.HomeZip, expected.SignedCirculoConsentForm, expected.CirculoConsentFormLink,
			expected.SignedStationMDConsentForm, expected.StationMDConsentFormLink, expected.CompletedGoSheet, expected.MarkedAsActive, expected.CreatedTimestamp,
			expected.LastModifiedTimestamp)

	s.mock.ExpectQuery(`SELECT * FROM "patient" WHERE patient_id = $1 ORDER BY "patient"."patient_id" LIMIT 1`).WithArgs("1").WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &PatientRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("1")

	s.Nil(err)
	s.Equal(expected, *result)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPatientRepo(t *testing.T) {
	suite.Run(t, new(PatientRepositorySuite))
}
