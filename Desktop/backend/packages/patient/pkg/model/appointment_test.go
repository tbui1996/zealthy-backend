package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AppointmentModelSuite struct {
	suite.Suite
}

func (s *AppointmentModelSuite) Test__IsNew() {
	model := &Appointment{}
	now := time.Now()

	s.True(model.IsNew())

	model.CreatedTimestamp = &now

	s.False(model.IsNew())
}

func (s *AppointmentModelSuite) Test__NewAppointment() {
	now := time.Now()
	model := NewAppointment(Appointment{
		AppointmentId:                 "1",
		PatientId:                     "2",
		AgencyProviderId:              "3",
		CirculatorDriverFullName:      "hey",
		AppointmentCreated:            &now,
		AppointmentScheduled:          "01/01/1990",
		AppointmentStatus:             "confirmed",
		AppointmentStatusChangedOn:    &now,
		AppointmentPurpose:            "nope",
		AppointmentOtherPurpose:       "",
		AppointmentNotes:              "",
		PatientDiastolicBloodPressure: 0,
		PatientSystolicBloodPressure:  0,
		PatientRespirationsPerMinute:  0,
		PatientPulseBeatsPerMinute:    0,
		PatientWeightLbs:              0,
		PatientChiefComplaint:         "rip",
		CreatedTimestamp:              &now,
		LastModifiedTimestamp:         &now,
	})

	s.Equal(model.AppointmentId, "1")
	s.Equal(model.PatientId, "2")
	s.Equal(model.AgencyProviderId, "3")
	s.Equal(model.AppointmentStatus, "confirmed")
}

func TestAppointmentModel(t *testing.T) {
	suite.Run(t, new(AppointmentModelSuite))
}
