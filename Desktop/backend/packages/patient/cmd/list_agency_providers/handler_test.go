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

type ListAgencyProviderSuite struct {
	suite.Suite
}

func (s *ListAgencyProviderSuite) Test__HandlesError() {
	mockRepo := new(mocks.AgencyProviderRepository)
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

func (s *ListAgencyProviderSuite) Test__Success() {

	expected := response.AgencyProviderResponse{
		AgencyProviderId:      "32f66c37-219b-46c7-82d5-f19c78b739ba",
		NationalProviderId:    "1",
		FirstName:             "bui",
		MiddleName:            "bui",
		LastName:              "bui",
		Suffix:                "Mr",
		BusinessName:          "circulo",
		BusinessTIN:           "what",
		BusinessAddress1:      "here",
		BusinessAddress2:      "",
		BusinessCity:          "dc",
		BusinessState:         "dc",
		BusinessZip:           "00000",
		CreatedTimestamp:      &time.Time{},
		LastModifiedTimestamp: &time.Time{},
	}
	mockRepo := new(mocks.AgencyProviderRepository)
	logger := zaptest.NewLogger(s.T())

	repoResults := []model.AgencyProvider{
		{
			AgencyProviderId:      "32f66c37-219b-46c7-82d5-f19c78b739ba",
			NationalProviderId:    "1",
			FirstName:             "bui",
			MiddleName:            "bui",
			LastName:              "bui",
			Suffix:                "Mr",
			BusinessName:          "circulo",
			BusinessTIN:           "what",
			BusinessAddress1:      "here",
			BusinessAddress2:      "",
			BusinessCity:          "dc",
			BusinessState:         "dc",
			BusinessZip:           "00000",
			CreatedTimestamp:      &time.Time{},
			LastModifiedTimestamp: &time.Time{},
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

func TestAgencyProviderHandler(t *testing.T) {
	suite.Run(t, new(ListAgencyProviderSuite))
}
