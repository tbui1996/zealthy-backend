package iface

import (
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
)

type AgencyProviderRepository interface {
	FindAll() (*[]model.AgencyProvider, *appointmenterror.AppointmentRepositoryError)
	Save(agencyProvider *model.AgencyProvider) *appointmenterror.AppointmentRepositoryError
	Find(agencyProviderId string) (*model.AgencyProvider, *appointmenterror.AppointmentRepositoryError)
}
