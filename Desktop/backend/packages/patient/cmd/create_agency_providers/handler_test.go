package main

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/patient/mocks"
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CreateAgencyProviderSuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sql.DB
	gormDb *gorm.DB
}

func (suite *CreateAgencyProviderSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	suite.db = db
	suite.mock = mock
	suite.gormDb = gdb
}

func (s *CreateAgencyProviderSuite) Test__HandlesSuccess() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateAgencyProviderRequest{
		NationalProviderId: "123",
		FirstName:          "thomas",
		DoddNumber:         "test",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessName:       "helloworld",
		BusinessCity:       "rock",
		BusinessState:      "md",
		BusinessZip:        "00002",
	}

	mockRepo.On("Save", mock.MatchedBy(func(agencyProvider *model.AgencyProvider) bool {
		return agencyProvider.NationalProviderId == "123" &&
			agencyProvider.FirstName == "thomas" &&
			agencyProvider.MiddleName == "" &&
			agencyProvider.DoddNumber == "test" &&
			agencyProvider.LastName == "bui" &&
			agencyProvider.Suffix == "" &&
			agencyProvider.BusinessAddress1 == "here" &&
			agencyProvider.BusinessAddress2 == "" &&
			agencyProvider.BusinessCity == "rock" &&
			agencyProvider.BusinessName == "helloworld" &&
			agencyProvider.BusinessState == "md" &&
			agencyProvider.BusinessZip == "00002"
	})).Return(nil)

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *CreateAgencyProviderSuite) Test__HandlesError() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}

	input := request.CreateAgencyProviderRequest{
		NationalProviderId: "123",
		FirstName:          "thomas",
		DoddNumber:         "test",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessName:       "helloworld",
		BusinessCity:       "rock",
		BusinessState:      "md",
		BusinessZip:        "00002",
	}

	mockRepo.On("Save", mock.MatchedBy(func(agencyProvider *model.AgencyProvider) bool {
		return agencyProvider.FirstName == "thomas" && agencyProvider.CreatedTimestamp == nil
	})).Return(appointmenterror.New("unknown error", appointmenterror.UNKNOWN))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func (s *CreateAgencyProviderSuite) Test__HandlesDuplicate() {
	mockRepo := new(mocks.AgencyProviderRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateAgencyProviderRequest{
		NationalProviderId: "123",
		FirstName:          "thomas",
		DoddNumber:         "test",
		MiddleName:         "",
		LastName:           "bui",
		Suffix:             "",
		BusinessTIN:        "",
		BusinessAddress1:   "here",
		BusinessAddress2:   "",
		BusinessName:       "helloworld",
		BusinessCity:       "rock",
		BusinessState:      "md",
		BusinessZip:        "00002",
	}

	mockRepo.On("Save", mock.MatchedBy(func(agencyProvider *model.AgencyProvider) bool {
		return agencyProvider.CreatedTimestamp == nil &&
			agencyProvider.FirstName == "thomas" &&
			agencyProvider.NationalProviderId == "123"
	})).Return(appointmenterror.New("dupe", appointmenterror.DODD_NUMBER_CONFLICT))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusConflict, result.StatusCode)
}

func TestCreateAgencyProviderHandler(t *testing.T) {
	suite.Run(t, new(CreateAgencyProviderSuite))
}
