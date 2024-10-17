package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.PatientRepository
	logger *zap.Logger
}

func Handler(input request.PatchPatientRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	patient, err := deps.repo.Find(input.PatientId)
	if err != nil {
		deps.logger.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
		return err.ToSonarApiGatewayResponse()
	}
	now := time.Now()
	newDateOfBirth, _ := time.Parse("2006-01-02", input.PatientDateOfBirth)
	patient.InsuranceId = input.InsuranceId
	patient.FirstName = input.PatientFirstName
	patient.MiddleName = input.PatientMiddleName
	patient.LastName = input.PatientLastName
	patient.EmailAddress = input.PatientEmailAddress
	patient.HomeAddress1 = input.PatientHomeAddress1
	patient.HomeAddress2 = input.PatientHomeAddress2
	patient.HomeCity = input.PatientHomeCity
	patient.HomeCounty = input.PatientHomeCounty
	patient.HomePhone = input.PatientHomePhone
	patient.HomeLivingArrangement = input.PatientHomeLivingArrangement
	patient.HomeZip = input.PatientHomeZip
	patient.DateOfBirth = newDateOfBirth
	patient.CirculoConsentFormLink = input.PatientCirculoConsentFormLink
	patient.InsuranceId = input.InsuranceId
	patient.LastModifiedTimestamp = &now
	patient.MarkedAsActive = input.PatientMarkedAsActive
	patient.CompletedGoSheet = input.PatientCompletedGoSheet
	patient.Suffix = input.PatientSuffix
	patient.SignedCirculoConsentForm = input.PatientSignedCirculoConsentForm
	patient.StationMDConsentFormLink = input.PatientStationMDConsentFormLink
	patient.HomeState = input.PatientHomeState
	patient.PreferredGender = input.PatientPreferredGender
	patient.PrimaryLanguage = input.PatientPrimaryLanguage
	patient.SignedStationMDConsentForm = input.PatientSignedStationMDConsentForm

	err = deps.repo.Save(patient)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
