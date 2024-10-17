package main

import (
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/patient/mocks"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CreateAppointmentSuite struct {
	suite.Suite
}

func (s *CreateAppointmentSuite) Test__HandlesSuccess() {
	mockRepo := new(mocks.AppointmentRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateAppointmentRequest{
		FirstName:                     "thomas",
		LastName:                      "bui",
		PatientId:                     "234",
		AgencyProviderId:              "1",
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "Confirmed",
		AppointmentPurpose:            "hello",
		PatientChiefComplaint:         "RIP",
		BusinessName:                  "agency",
		CirculatorDriverFullName:      "me",
		PatientSystolicBloodPressure:  0,
		PatientDiastolicBloodPressure: 0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
	}

	mockRepo.On("Save", mock.MatchedBy(func(appointment *model.Appointment) bool {
		return appointment.PatientId == "234" &&
			appointment.AppointmentStatus == "Confirmed" &&
			appointment.PatientChiefComplaint == "RIP" &&
			appointment.AgencyProviderId == "1"
	})).Return(nil)

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func TestCreateAppointmentHandler(t *testing.T) {
	suite.Run(t, new(CreateAppointmentSuite))
}
