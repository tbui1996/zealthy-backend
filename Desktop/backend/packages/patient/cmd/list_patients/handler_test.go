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

type ListPatientSuite struct {
	suite.Suite
}

func (s *ListPatientSuite) Test__HandlesError() {
	mockRepo := new(mocks.PatientRepository)
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

func (s *ListPatientSuite) Test__Success() {

	expected := response.PatientResponse{
		PatientId:                  "2d19bdba-1f9a-444d-999f-65172f581cfc",
		InsuranceId:                "1237918237",
		FirstName:                  "thomas",
		MiddleName:                 "bui",
		LastName:                   "bui",
		Suffix:                     "MR",
		DateOfBirth:                &time.Time{},
		PrimaryLanguage:            "english",
		PreferredGender:            "male",
		EmailAddress:               "",
		HomePhone:                  "0000000000",
		HomeLivingArrangement:      "yes",
		HomeAddress1:               "oh",
		HomeAddress2:               "io",
		HomeCity:                   "columbus",
		HomeCounty:                 "moco",
		HomeState:                  "md",
		HomeZip:                    "00000",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "www.google.com",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "www.bing.com",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
		CreatedTimestamp:           &time.Time{},
		LastModifiedTimestamp:      &time.Time{},
	}
	mockRepo := new(mocks.PatientRepository)
	logger := zaptest.NewLogger(s.T())

	repoResults := []model.Patient{
		{
			PatientId:                  "2d19bdba-1f9a-444d-999f-65172f581cfc",
			InsuranceId:                "1237918237",
			FirstName:                  "thomas",
			MiddleName:                 "bui",
			LastName:                   "bui",
			Suffix:                     "MR",
			DateOfBirth:                time.Time{},
			PrimaryLanguage:            "english",
			PreferredGender:            "male",
			EmailAddress:               "",
			HomePhone:                  "0000000000",
			HomeLivingArrangement:      "yes",
			HomeAddress1:               "oh",
			HomeAddress2:               "io",
			HomeCity:                   "columbus",
			HomeCounty:                 "moco",
			HomeState:                  "md",
			HomeZip:                    "00000",
			SignedCirculoConsentForm:   false,
			CirculoConsentFormLink:     "www.google.com",
			SignedStationMDConsentForm: false,
			StationMDConsentFormLink:   "www.bing.com",
			CompletedGoSheet:           false,
			MarkedAsActive:             false,
			CreatedTimestamp:           &time.Time{},
			LastModifiedTimestamp:      &time.Time{},
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

func TestListFlagsHandler(t *testing.T) {
	suite.Run(t, new(ListPatientSuite))
}
