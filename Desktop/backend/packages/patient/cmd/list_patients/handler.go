package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/response"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.PatientRepository
	logger *zap.Logger
}

func Handler(deps HandlerDeps) (*[]response.PatientResponse, error) {
	patients, err := deps.repo.FindAll()

	if err != nil {
		deps.logger.Error(err.Error())
		return nil, fmt.Errorf("unable to find records in database")
	}

	responses := make([]response.PatientResponse, len(*patients))

	for i := range *patients {
		patient := (*patients)[i]
		responses[i] = response.PatientResponse{
			PatientId:                  patient.PatientId,
			InsuranceId:                patient.InsuranceId,
			FirstName:                  patient.FirstName,
			MiddleName:                 patient.MiddleName,
			LastName:                   patient.LastName,
			Suffix:                     patient.Suffix,
			DateOfBirth:                &patient.DateOfBirth,
			EmailAddress:               patient.EmailAddress,
			PrimaryLanguage:            patient.PrimaryLanguage,
			PreferredGender:            patient.PreferredGender,
			HomePhone:                  patient.HomePhone,
			HomeLivingArrangement:      patient.HomeLivingArrangement,
			HomeAddress1:               patient.HomeAddress1,
			HomeAddress2:               patient.HomeAddress2,
			HomeCity:                   patient.HomeCity,
			HomeCounty:                 patient.HomeCounty,
			HomeState:                  patient.HomeState,
			HomeZip:                    patient.HomeZip,
			SignedCirculoConsentForm:   patient.SignedCirculoConsentForm,
			CirculoConsentFormLink:     patient.CirculoConsentFormLink,
			SignedStationMDConsentForm: patient.SignedStationMDConsentForm,
			StationMDConsentFormLink:   patient.StationMDConsentFormLink,
			CompletedGoSheet:           patient.CompletedGoSheet,
			MarkedAsActive:             patient.MarkedAsActive,
			CreatedTimestamp:           patient.CreatedTimestamp,
			LastModifiedTimestamp:      patient.LastModifiedTimestamp,
		}
	}
	return &responses, nil
}
