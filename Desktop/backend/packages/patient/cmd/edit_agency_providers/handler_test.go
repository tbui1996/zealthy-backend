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

type EditAgencyProviderSuite struct {
	suite.Suite
}

func (s *EditAgencyProviderSuite) Test__HandleEditSuccess() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}

	input := request.EditAgencyProviderRequest{
		AgencyProviderId:   "",
		NationalProviderId: "",
		DoddNumber:         "",
		FirstName:          "",
		MiddleName:         "",
		LastName:           "",
		Suffix:             "",
		BusinessName:       "",
		BusinessTIN:        "",
		BusinessAddress1:   "",
		BusinessAddress2:   "",
		BusinessCity:       "",
		BusinessState:      "",
		BusinessZip:        "",
	}
	agencyProvider := &model.AgencyProvider{
		AgencyProviderId:      "",
		NationalProviderId:    "",
		FirstName:             "hola",
		DoddNumber:            "123",
		MiddleName:            "me",
		LastName:              "llamo",
		Suffix:                "MR",
		BusinessName:          "",
		BusinessTIN:           "",
		BusinessAddress1:      "",
		BusinessAddress2:      "",
		BusinessCity:          "",
		BusinessState:         "",
		BusinessZip:           "",
		CreatedTimestamp:      &time.Time{},
		LastModifiedTimestamp: &time.Time{},
	}

	mockRepo.On("Find", "").Return(agencyProvider, nil)
	mockRepo.On("Save", agencyProvider).Return(nil)

	result, err := Handler(input, deps)
	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(agencyProvider *model.AgencyProvider) bool {
		return agencyProvider.FirstName == input.FirstName && agencyProvider.MiddleName == input.MiddleName
	}))

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *EditAgencyProviderSuite) Test__HandleEditFailFind() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}

	input := request.EditAgencyProviderRequest{
		DoddNumber:         "2",
		AgencyProviderId:   "3",
		NationalProviderId: "4",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessCity:       "columbus",
		BusinessState:      "ohio",
		BusinessZip:        "80000",
	}

	mockRepo.On("Find", "3").Return(nil, appointmenterror.New("FAKE ERROR, IGNORE", appointmenterror.UNKNOWN))

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", "3")
	mockRepo.AssertNotCalled(s.T(), "Save", mock.Anything)
	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func (s *EditAgencyProviderSuite) Test__HandleEditFailSave() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}

	input := request.EditAgencyProviderRequest{
		DoddNumber:         "2",
		AgencyProviderId:   "3",
		NationalProviderId: "4",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessName:       "circulo",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessCity:       "columbus",
		BusinessState:      "ohio",
		BusinessZip:        "80000",
	}

	agencyProvider := &model.AgencyProvider{
		AgencyProviderId:   "3",
		DoddNumber:         "2",
		NationalProviderId: "4",
		FirstName:          "thomas",
		MiddleName:         "",
		LastName:           "bu",
		Suffix:             "",
		BusinessName:       "circul",
		BusinessTIN:        "",
		BusinessAddress1:   "her",
		BusinessAddress2:   "",
		BusinessCity:       "columbu",
		BusinessState:      "ohi",
		BusinessZip:        "80000",
	}

	mockRepo.On("Find", "3").Return(agencyProvider, nil)
	mockRepo.On("Save", agencyProvider).Return(appointmenterror.New("FAKE ERROR, IGNORE", appointmenterror.UNKNOWN))

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", "3")
	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(agencyProvider *model.AgencyProvider) bool {
		return agencyProvider.BusinessName == input.BusinessName && agencyProvider.FirstName == input.FirstName
	}))
	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func TestEditAgencyProviderHandler(t *testing.T) {
	suite.Run(t, new(EditAgencyProviderSuite))
}
