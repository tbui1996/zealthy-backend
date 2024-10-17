package main

import (
	"encoding/json"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/patients/iface"
	"go.uber.org/zap"
)

type PatientsGetRequest struct {
	UserId string
	Repo   iface.PatientRepository
	Logger *zap.Logger
}

type PatientResponse struct {
	Patients []model.Patient `json:"patients"`
}

func Handler(req PatientsGetRequest) ([]byte, error) {
	patients, err := req.Repo.FindAll(model.Patient{ProviderId: req.UserId})

	if err != nil {
		req.Logger.Error("unable to get patients for provider")
		return nil, err
	}

	patientBytes, err := json.Marshal(PatientResponse{Patients: patients})
	if err != nil {
		req.Logger.Error("unable to marshal patient object")
		return nil, err
	}

	return patientBytes, nil
}
