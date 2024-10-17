package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/patient/mocks"
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CreatePatientSuite struct {
	suite.Suite
}

func (s *CreatePatientSuite) Test__HandlesDuplicate() {
	mockRepo := new(mocks.PatientRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreatePatientRequest{
		InsuranceId:                     "123",
		PatientFirstName:                "thomas",
		PatientMiddleName:               "",
		PatientLastName:                 "bui",
		PatientSuffix:                   "",
		PatientDateOfBirth:              "0001-01-01",
		PatientPrimaryLanguage:          "",
		PatientPreferredGender:          "",
		PatientEmailAddress:             "thomas@circulohealth.com",
		PatientHomePhone:                "",
		PatientHomeLivingArrangement:    "",
		PatientHomeAddress1:             "here",
		PatientHomeAddress2:             "",
		PatientHomeCity:                 "rock",
		PatientHomeCounty:               "moco",
		PatientHomeState:                "md",
		PatientHomeZip:                  "00000",
		PatientCirculoConsentFormLink:   "",
		PatientStationMDConsentFormLink: "",
	}

	mockRepo.On("Save", mock.MatchedBy(func(patient *model.Patient) bool {
		return patient.InsuranceId == "123" &&
			patient.FirstName == "thomas" &&
			patient.LastName == "bui" &&
			patient.DateOfBirth == time.Time{} &&
			patient.EmailAddress == "thomas@circulohealth.com" &&
			patient.HomeAddress1 == "here" &&
			patient.HomeCity == "rock" &&
			patient.HomeCounty == "moco" &&
			patient.HomeState == "md" &&
			patient.HomeZip == "00000"
	})).Return(appointmenterror.New("dupe", appointmenterror.INSURANCE_ID_CONFLICT))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusConflict, result.StatusCode)
}

func (s *CreatePatientSuite) Test__HandlesSuccess() {
	mockRepo := new(mocks.PatientRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreatePatientRequest{
		InsuranceId:                     "123",
		PatientFirstName:                "thomas",
		PatientMiddleName:               "",
		PatientLastName:                 "bui",
		PatientSuffix:                   "",
		PatientDateOfBirth:              "0001-01-01",
		PatientPrimaryLanguage:          "",
		PatientPreferredGender:          "",
		PatientEmailAddress:             "thomas@circulohealth.com",
		PatientHomePhone:                "",
		PatientHomeLivingArrangement:    "",
		PatientHomeAddress1:             "here",
		PatientHomeAddress2:             "",
		PatientHomeCity:                 "rock",
		PatientHomeCounty:               "moco",
		PatientHomeState:                "md",
		PatientHomeZip:                  "00002",
		PatientCirculoConsentFormLink:   "",
		PatientStationMDConsentFormLink: "",
	}

	mockRepo.On("Save", mock.MatchedBy(func(patient *model.Patient) bool {
		return patient.InsuranceId == "123" &&
			patient.FirstName == "thomas" &&
			patient.LastName == "bui" &&
			patient.DateOfBirth == time.Time{} &&
			patient.EmailAddress == "thomas@circulohealth.com" &&
			patient.HomeAddress1 == "here" &&
			patient.HomeCity == "rock" &&
			patient.HomeCounty == "moco" &&
			patient.HomeState == "md" &&
			patient.HomeZip == "00002"
	})).Return(nil)

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *CreatePatientSuite) Test__HandlesError() {
	mockRepo := new(mocks.PatientRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreatePatientRequest{
		InsuranceId:                     "",
		PatientFirstName:                "thomas",
		PatientMiddleName:               "",
		PatientLastName:                 "bui",
		PatientSuffix:                   "",
		PatientDateOfBirth:              "0001-01-01",
		PatientPrimaryLanguage:          "",
		PatientPreferredGender:          "",
		PatientEmailAddress:             "thomas@circulohealth.com",
		PatientHomePhone:                "",
		PatientHomeLivingArrangement:    "",
		PatientHomeAddress1:             "here",
		PatientHomeAddress2:             "",
		PatientHomeCity:                 "rock",
		PatientHomeCounty:               "moco",
		PatientHomeState:                "md",
		PatientHomeZip:                  "00002",
		PatientCirculoConsentFormLink:   "",
		PatientStationMDConsentFormLink: "",
	}

	mockRepo.On("Save", mock.MatchedBy(func(patient *model.Patient) bool {
		return patient.InsuranceId == "" &&
			patient.FirstName == "thomas" &&
			patient.LastName == "bui" &&
			patient.DateOfBirth == time.Time{} &&
			patient.EmailAddress == "thomas@circulohealth.com" &&
			patient.HomeAddress1 == "here" &&
			patient.HomeCity == "rock" &&
			patient.HomeCounty == "moco" &&
			patient.HomeState == "md" &&
			patient.HomeZip == "00002" &&
			patient.CreatedTimestamp == nil
	})).Return(appointmenterror.New("unknown error", appointmenterror.UNKNOWN))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusBadRequest, result.StatusCode)
}

func TestCreatePatientHandler(t *testing.T) {
	suite.Run(t, new(CreatePatientSuite))
}
