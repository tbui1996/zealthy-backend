package data

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AppointmentRepositorySuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
}

func (suite *AppointmentRepositorySuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (s *AppointmentRepositorySuite) Test__FindAllSuccess() {
	now := time.Now()
	expected := model.JoinResult{
		AppointmentId:                 "1",
		PatientId:                     "1",
		AgencyProviderId:              "1",
		CirculatorDriverFullName:      "hi",
		AppointmentCreated:            &now,
		AppointmentScheduled:          &now,
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "yes",
		AppointmentOtherPurpose:       "oops",
		AppointmentNotes:              "notes",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "WHY",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
		FirstName:                     "thomas",
		MiddleName:                    "h",
		LastName:                      "bui",
		ProviderFullName:              "charles",
		Suffix:                        "MR",
		DateOfBirth:                   time.Time{},
		PrimaryLanguage:               "eng",
		PreferredGender:               "M",
		EmailAddress:                  "thomas@circulohealth.com",
		HomeAddress1:                  "here",
		HomeAddress2:                  "",
		HomeCity:                      "columbus",
		HomeState:                     "oh",
		HomeZip:                       "00000",
		SignedCirculoConsentForm:      false,
		CirculoConsentFormLink:        "",
		SignedStationMDConsentForm:    false,
		StationMDConsentFormLink:      "",
		CompletedGoSheet:              false,
		MarkedAsActive:                false,
		NationalProviderId:            "123",
		BusinessName:                  "circulo",
		BusinessTIN:                   "",
		BusinessAddress1:              "yes",
		BusinessAddress2:              "",
		BusinessCity:                  "columbus",
		BusinessState:                 "oh",
		BusinessZip:                   "00002",
		HomePhone:                     "",
		HomeLivingArrangement:         "",
		HomeCounty:                    "",
		InsuranceId:                   "234",
	}
	rows := sqlmock.
		NewRows([]string{"appointment_id", "patient_id", "agency_provider_id",
			"circulator_driver_fullname", "appointment_created", "appointment_scheduled",
			"appointment_status", "appointment_status_changed_on", "appointment_purpose",
			"appointment_other_purpose", "appointment_notes", "patient_diastolic_blood_pressure",
			"patient_systolic_blood_pressure", "patient_respirations_per_minute", "patient_pulse_beats_per_minute",
			"patient_weight_lbs", "patient_chief_complaint", "created_timestamp",
			"last_modified_timestamp", "first_name", "middle_name",
			"last_name", "provider_fullname", "suffix",
			"date_of_birth", "primary_language", "preferred_gender",
			"email_address", "home_address_1", "home_address_2",
			"home_city", "home_state", "home_zip",
			"signed_circulo_consent_form", "circulo_consent_form_link", "signed_stationmd_consent_form",
			"stationmd_consent_form_link", "completed_go_sheet", "marked_as_active",
			"national_provider_id", "business_name", "business_tin",
			"business_address_1", "business_address_2", "business_city",
			"business_state", "business_zip", "home_phone",
			"home_living_arrangement", "home_county", "insurance_id"}).
		AddRow("1", "1", "1",
			"hi", now, now,
			"confirmed", now, "yes",
			"oops", "notes", 0,
			0, 0, 0,
			0, "WHY", now,
			now, "thomas", "h",
			"bui", "charles", "MR",
			time.Time{}, "eng", "M",
			"thomas@circulohealth.com", "here", "",
			"columbus", "oh", "00000",
			false, "", false, "",
			false, false, "123",
			"circulo", "", "yes",
			"", "columbus", "oh",
			"00002", "", "",
			"", "234")
	s.mock.ExpectQuery(`SELECT a.appointment_id, a.appointment_status, a.appointment_purpose, a.appointment_notes, a.appointment_other_purpose,a.appointment_created, a.appointment_scheduled, a.appointment_status_changed_on, a.circulator_driver_fullname, a.patient_diastolic_blood_pressure, a.patient_systolic_blood_pressure, a.patient_respirations_per_minute, a.patient_pulse_beats_per_minute, a.patient_weight_lbs, a.patient_chief_complaint, b.*, c.agency_provider_id, c.national_provider_id, CONCAT_WS(' ', c.first_name, c.middle_name, c.last_name, ' ', c.suffix) as provider_fullname, c.business_name, c.business_tin, c.business_address_1, c.business_address_2, c.business_city, c.business_state, c.business_zip FROM appointments as a JOIN patient as b on a.patient_id = b.patient_id JOIN agency_provider as c on a.agency_provider_id = c.agency_provider_id ORDER BY a.last_modified_timestamp desc`).
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
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

func (s *AppointmentRepositorySuite) Test__FindSuccess() {
	now := time.Now()
	expected := model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "123",
		AgencyProviderId:              "123",
		CirculatorDriverFullName:      "",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "here",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "yes",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
	}
	rows := sqlmock.
		NewRows([]string{"appointment_id", "patient_id", "agency_provider_id",
			"circulator_driver_fullname", "appointment_created", "appointment_scheduled",
			"appointment_status", "appointment_status_changed_on", "appointment_purpose",
			"appointment_other_purpose", "appointment_notes", "patient_diastolic_blood_pressure",
			"patient_systolic_blood_pressure", "patient_respirations_per_minute", "patient_pulse_beats_per_minute",
			"patient_weight_lbs", "patient_chief_complaint", "created_timestamp", "last_modified_timestamp"}).
		AddRow("123", "123", "123", "", now, "01/01/1990", "confirmed", now, "here", "", "", 0, 0, 0, 0, 0, "yes", now, now)
	s.mock.ExpectQuery(`SELECT * FROM "appointments" WHERE appointment_id = $1 ORDER BY "appointments"."appointment_id" LIMIT 1`).
		WithArgs("123").
		WillReturnRows(rows)

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("123")

	s.Nil(err)
	s.Equal(expected, *result)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

}

func (s *AppointmentRepositorySuite) Test__FindFailNotFound() {
	s.mock.ExpectQuery(`SELECT * FROM "appointments" WHERE appointment_id = $1 ORDER BY "appointments"."appointment_id" LIMIT 1`).
		WithArgs("1").
		WillReturnError(gorm.ErrRecordNotFound)

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("1")
	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AppointmentRepositorySuite) Test__FindFail() {
	s.mock.ExpectQuery(`SELECT * FROM "appointments" WHERE appointment_id = $1 ORDER BY "appointments"."appointment_id" LIMIT 1`).
		WithArgs("1").
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	result, err := repo.Find("1")
	s.NotNil(err)
	s.Nil(result)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AppointmentRepositorySuite) Test__DeleteSuccess() {
	now := time.Now()
	expected := model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "123",
		AgencyProviderId:              "123",
		CirculatorDriverFullName:      "",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "here",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "yes",
		CreatedTimestamp:              &time.Time{},
		LastModifiedTimestamp:         &time.Time{},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "appointments" WHERE "appointments"."appointment_id" = $1`).
		WithArgs("123").
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Delete(&expected)
	s.Nil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__DeleteFail() {
	now := time.Now()
	expected := model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "123",
		AgencyProviderId:              "123",
		CirculatorDriverFullName:      "",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "here",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "yes",
		CreatedTimestamp:              &time.Time{},
		LastModifiedTimestamp:         &time.Time{},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "appointments" WHERE "appointments"."appointment_id" = $1`).
		WithArgs("123").
		WillReturnError(errors.New("FAKE ERROR, IGNORE"))
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Delete(&expected)
	s.NotNil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *AgencyProviderRepositorySuite) Test__DeleteFailNotFound() {
	now := time.Now()
	expected := model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "123",
		AgencyProviderId:              "123",
		CirculatorDriverFullName:      "",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "here",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "yes",
		CreatedTimestamp:              &time.Time{},
		LastModifiedTimestamp:         &time.Time{},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "appointments" WHERE "appointments"."appointment_id" = $1`).
		WithArgs("123").
		WillReturnError(gorm.ErrRecordNotFound)
	s.mock.ExpectRollback()

	logger := zaptest.NewLogger(s.T())

	repo := &AppointmentRepository{
		DB:     s.gormDb,
		Logger: logger,
	}

	err := repo.Delete(&expected)
	s.NotNil(err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAppointmentRepo(t *testing.T) {
	suite.Run(t, new(AppointmentRepositorySuite))
}
