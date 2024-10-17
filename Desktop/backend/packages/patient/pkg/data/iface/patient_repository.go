package iface

import (
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
)

type PatientRepository interface {
	FindAll() (*[]model.Patient, *appointmenterror.AppointmentRepositoryError)
	Save(patient *model.Patient) *appointmenterror.AppointmentRepositoryError
	Find(patientId string) (*model.Patient, *appointmenterror.AppointmentRepositoryError)
}
