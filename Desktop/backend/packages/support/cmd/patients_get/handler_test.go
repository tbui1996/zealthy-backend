package main

import (
	"errors"
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/patients/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type PatientsGetSuite struct {
	suite.Suite
	req  PatientsGetRequest
	repo *mocks.PatientRepository
}

func (suite *PatientsGetSuite) SetupTest() {
	repo := new(mocks.PatientRepository)

	suite.repo = repo
	suite.req = PatientsGetRequest{
		UserId: "1",
		Repo:   repo,
		Logger: zaptest.NewLogger(suite.T()),
	}
}

func (suite *PatientsGetSuite) TestPatientsGet_Success() {
	p := model.Patient{
		ProviderId:  "1",
		ID:          1,
		LastName:    "User",
		Name:        "Test",
		Address:     "123 Address Street",
		InsuranceID: "123456789",
		Birthday:    time.Now(),
	}

	suite.repo.On("FindAll", mock.Anything).Return([]model.Patient{p}, nil)

	out, err := Handler(suite.req)

	suite.NotNil(out)
	suite.NoError(err)
}

func (suite *PatientsGetSuite) TestPatientsGet_Fail() {
	suite.repo.On("FindAll", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	out, err := Handler(suite.req)

	suite.Nil(out)
	suite.Error(err)
}

func TestPatientsGet(t *testing.T) {
	suite.Run(t, new(PatientsGetSuite))
}
