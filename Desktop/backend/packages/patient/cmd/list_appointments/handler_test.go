package main

import (
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/patient/mocks"
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/response"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ListAppointmentSuite struct {
	suite.Suite
}

func (s *ListAppointmentSuite) Test__HandlesError() {
	mockRepo := new(mocks.AppointmentRepository)
	logger := zaptest.NewLogger(s.T())

	mockRepo.On("FindAll").Return(nil, appointmenterror.New("unknown error", appointmenterror.UNKNOWN))

	deps := HandlerDeps{
		repo:   mockRepo,
		logger: logger,
	}

	results, err := Handler(deps)

	s.Nil(results)
	s.NotNil(err)
}

func (s *ListAppointmentSuite) Test__Success() {
	now := time.Now()
	expected := response.AppointmentResponse{
		AppointmentId:                 "291400f3-2820-42c3-89c5-f68cb3e74f24",
		PatientId:                     "2d19bdba-1f9a-444d-999f-65172f581cfc",
		AgencyProviderId:              "32f66c37-219b-46c7-82d5-f19c78b739ba",
		CirculatorDriverFullName:      "my man",
		AppointmentCreated:            &now,
		AppointmentScheduled:          &now,
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "oopsies",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
		FirstName:                     "bui",
		MiddleName:                    "bui",
		LastName:                      "bui",
		ProviderFullName:              "d d d",
		Suffix:                        "",
		DateOfBirth:                   &now,
		PrimaryLanguage:               "",
		PreferredGender:               "",
		EmailAddress:                  "",
		HomeAddress1:                  "",
		HomeAddress2:                  "",
		HomeCity:                      "",
		HomeState:                     "",
		HomeZip:                       "",
		SignedCirculoConsentForm:      false,
		CirculoConsentFormLink:        "",
		SignedStationMDConsentForm:    false,
		StationMDConsentFormLink:      "",
		CompletedGoSheet:              false,
		MarkedAsActive:                false,
		NationalProviderId:            "",
		BusinessName:                  "",
		BusinessTIN:                   "",
		BusinessAddress1:              "",
		BusinessAddress2:              "",
		BusinessCity:                  "",
		BusinessState:                 "",
		BusinessZip:                   "",
		HomePhone:                     "",
		HomeLivingArrangement:         "",
		HomeCounty:                    "",
		InsuranceId:                   "",
	}
	mockRepo := new(mocks.AppointmentRepository)
	logger := zaptest.NewLogger(s.T())

	repoResults := []model.JoinResult{
		{
			AppointmentId:                 "291400f3-2820-42c3-89c5-f68cb3e74f24",
			PatientId:                     "2d19bdba-1f9a-444d-999f-65172f581cfc",
			AgencyProviderId:              "32f66c37-219b-46c7-82d5-f19c78b739ba",
			CirculatorDriverFullName:      "my man",
			AppointmentCreated:            &now,
			AppointmentScheduled:          &now,
			AppointmentStatus:             "confirmed",
			AppointmentStatusChangedOn:    &now,
			AppointmentPurpose:            "",
			AppointmentOtherPurpose:       "",
			AppointmentNotes:              "",
			PatientDiastolicBloodPressure: 0,
			PatientSystolicBloodPressure:  0,
			PatientRespirationsPerMinute:  0,
			PatientPulseBeatsPerMinute:    0,
			PatientWeightLbs:              0,
			PatientChiefComplaint:         "oopsies",
			CreatedTimestamp:              &now,
			LastModifiedTimestamp:         &now,
			FirstName:                     "bui",
			MiddleName:                    "bui",
			LastName:                      "bui",
			ProviderFullName:              "d d d",
			Suffix:                        "",
			DateOfBirth:                   now,
			PrimaryLanguage:               "",
			PreferredGender:               "",
			EmailAddress:                  "",
			HomeAddress1:                  "",
			HomeAddress2:                  "",
			HomeCity:                      "",
			HomeState:                     "",
			HomeZip:                       "",
			SignedCirculoConsentForm:      false,
			CirculoConsentFormLink:        "",
			SignedStationMDConsentForm:    false,
			StationMDConsentFormLink:      "",
			CompletedGoSheet:              false,
			MarkedAsActive:                false,
			NationalProviderId:            "",
			BusinessName:                  "",
			BusinessTIN:                   "",
			BusinessAddress1:              "",
			BusinessAddress2:              "",
			BusinessCity:                  "",
			BusinessState:                 "",
			BusinessZip:                   "",
			HomePhone:                     "",
			HomeLivingArrangement:         "",
			HomeCounty:                    "",
			InsuranceId:                   "",
		},
	}

	mockRepo.On("FindAll").Return(&repoResults, nil)

	deps := HandlerDeps{
		repo:   mockRepo,
		logger: logger,
	}

	results, err := Handler(deps)

	s.Equal(expected, (*results)[0])
	s.Nil(err)
}

func TestAppointmentHandler(t *testing.T) {
	suite.Run(t, new(ListAppointmentSuite))
}
