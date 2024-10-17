package iface

import (
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
)

type AppointmentRepository interface {
	FindAll() (*[]model.JoinResult, *appointmenterror.AppointmentRepositoryError)
	Save(appointment *model.Appointment) *appointmenterror.AppointmentRepositoryError
	Find(appointmentId string) (*model.Appointment, *appointmenterror.AppointmentRepositoryError)
	Delete(appointment *model.Appointment) *appointmenterror.AppointmentRepositoryError
}
