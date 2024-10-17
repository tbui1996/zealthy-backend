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

type DeleteAppointmentSuite struct {
	suite.Suite
}

func (s *DeleteAppointmentSuite) Test__HandleDeleteSuccess() {
	now := time.Now()
	mockRepo := new(mocks.AppointmentRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteAppointmentRequest{
		AppointmentId: "123",
	}
	appointment := &model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "12",
		AgencyProviderId:              "12",
		CirculatorDriverFullName:      "23",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
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
		PatientChiefComplaint:         "",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
	}

	mockRepo.On("Find", "123").Return(appointment, nil)
	mockRepo.On("Delete", appointment).Return(nil)

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Delete", mock.Anything)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *DeleteAppointmentSuite) Test__HandleDeleteFailId() {
	mockRepo := new(mocks.AppointmentRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteAppointmentRequest{
		AppointmentId: "",
	}

	result, _ := Handler(input, deps)

	mockRepo.AssertNotCalled(s.T(), "Find", mock.Anything)
	mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything)
	s.Equal(http.StatusBadRequest, result.StatusCode)
}

func (s *DeleteAppointmentSuite) Test__HandlerDeleteFailFind() {
	now := time.Now()
	mockRepo := new(mocks.AppointmentRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteAppointmentRequest{
		AppointmentId: "123",
	}

	appointment := &model.Appointment{
		AppointmentId:                 "123",
		PatientId:                     "12",
		AgencyProviderId:              "12",
		CirculatorDriverFullName:      "23",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
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
		PatientChiefComplaint:         "",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
	}

	mockRepo.On("Find", "123").Return(nil, appointmenterror.New("Fake error, ignore", appointmenterror.NOT_FOUND))
	mockRepo.On("Delete", appointment).Return(nil)

	result, _ := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", mock.Anything)
	mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything)

	s.Equal(http.StatusNotFound, result.StatusCode)
}

func TestDeleteAppointmentHandler(t *testing.T) {
	suite.Run(t, new(DeleteAppointmentSuite))
}
