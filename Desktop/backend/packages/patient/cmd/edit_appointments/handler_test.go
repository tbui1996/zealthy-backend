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

type EditAppointmentSuite struct {
	suite.Suite
}

func (s *EditAppointmentSuite) Test__HandleEditSuccess() {
	mockRepo := new(mocks.AppointmentRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	now := time.Now()

	input := request.EditAppointmentRequest{
		AppointmentId:                 "12",
		FirstName:                     "",
		LastName:                      "",
		PatientId:                     "",
		AgencyProviderId:              "",
		AppointmentScheduled:          "",
		AppointmentStatus:             "",
		AppointmentPurpose:            "",
		PatientChiefComplaint:         "",
		BusinessName:                  "",
		CirculatorDriverFullName:      "",
		PatientSystolicBloodPressure:  0,
		PatientDiastolicBloodPressure: 0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
	}

	appointment := &model.Appointment{
		AppointmentId:                 "12",
		PatientId:                     "23",
		AgencyProviderId:              "13",
		CirculatorDriverFullName:      "circulator",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "ruhroh",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
	}

	mockRepo.On("Find", "12").Return(appointment, nil)
	mockRepo.On("Save", appointment).Return(nil)

	result, err := Handler(input, deps)
	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(appointment *model.Appointment) bool {
		return appointment.AppointmentId == input.AppointmentId
	}))

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func TestEditAppointmentHandler(t *testing.T) {
	suite.Run(t, new(EditAppointmentSuite))
}
