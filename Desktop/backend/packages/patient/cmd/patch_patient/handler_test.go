package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/patient/mocks"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type PatchPatientSuite struct {
	suite.Suite
}

func (s *PatchPatientSuite) Test__HandlePatchSuccess() {
	mockRepo := new(mocks.PatientRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}

	input := request.PatchPatientRequest{
		InsuranceId:                     "0",
		PatientFirstName:                "hola",
		PatientMiddleName:               "me",
		PatientLastName:                 "llamo",
		PatientSuffix:                   "MR",
		PatientDateOfBirth:              "1996-12-28",
		PatientPrimaryLanguage:          "",
		PatientPreferredGender:          "",
		PatientEmailAddress:             "thomas",
		PatientHomePhone:                "",
		PatientHomeLivingArrangement:    "",
		PatientHomeAddress1:             "yes",
		PatientHomeAddress2:             "",
		PatientHomeCity:                 "potomac",
		PatientHomeCounty:               "no",
		PatientHomeState:                "dc",
		PatientHomeZip:                  "12345",
		PatientCirculoConsentFormLink:   "",
		PatientStationMDConsentFormLink: "",
	}
	patient := &model.Patient{
		PatientId:                  "123",
		InsuranceId:                "",
		FirstName:                  "hola",
		MiddleName:                 "me",
		LastName:                   "llamo",
		Suffix:                     "MR",
		DateOfBirth:                time.Time{},
		PrimaryLanguage:            "",
		PreferredGender:            "",
		EmailAddress:               "uhoh",
		HomePhone:                  "",
		HomeLivingArrangement:      "",
		HomeAddress1:               "yes",
		HomeAddress2:               "",
		HomeCity:                   "up",
		HomeCounty:                 "dog",
		HomeState:                  "jpws",
		HomeZip:                    "2879",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
		CreatedTimestamp:           &time.Time{},
		LastModifiedTimestamp:      &time.Time{},
	}

	mockRepo.On("Find", "").Return(patient, nil)
	mockRepo.On("Save", patient).Return(nil)

	result, err := Handler(input, deps)
	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(patient *model.Patient) bool {
		return patient.FirstName == input.PatientFirstName && patient.MiddleName == input.PatientMiddleName
	}))

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func TestPatchPatientHandler(t *testing.T) {
	suite.Run(t, new(PatchPatientSuite))
}
