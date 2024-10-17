package iface

import "github.com/circulohealth/sonar-backend/packages/support/pkg/model"

type PatientRepository interface {
	FindAll(filter interface{}) ([]model.Patient, error)
	Find(filter interface{}) (*model.Patient, error)
}
