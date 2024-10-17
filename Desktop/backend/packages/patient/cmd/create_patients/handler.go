package main

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.PatientRepository
	logger *zap.Logger
}

func Handler(input request.CreatePatientRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {

	if input.InsuranceId == "" || input.PatientDateOfBirth == "" || input.PatientFirstName == "" || input.PatientLastName == "" || input.PatientEmailAddress == "" || input.PatientHomeAddress1 == "" || input.PatientHomeCity == "" || input.PatientHomeCounty == "" || input.PatientHomeState == "" || input.PatientHomeZip == "" {
		deps.logger.Error("Insurance ID, First/Last Name, Date of Birth, Email Address, Home address 1, City, County, State, Zip Code are all required")
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "Insurance ID, First/Last Name, Date of Birth, Email Address, Home address 1, City, County, State, Zip Code are all required")
	}

	myDate, _ := time.Parse("2006-01-02", input.PatientDateOfBirth)

	patient := &model.Patient{
		InsuranceId:              input.InsuranceId,
		FirstName:                input.PatientFirstName,
		MiddleName:               input.PatientMiddleName,
		LastName:                 input.PatientLastName,
		Suffix:                   input.PatientSuffix,
		DateOfBirth:              myDate,
		PrimaryLanguage:          input.PatientPrimaryLanguage,
		PreferredGender:          input.PatientPreferredGender,
		EmailAddress:             input.PatientEmailAddress,
		HomePhone:                input.PatientHomePhone,
		HomeLivingArrangement:    input.PatientHomeLivingArrangement,
		HomeAddress1:             input.PatientHomeAddress1,
		HomeAddress2:             input.PatientHomeAddress2,
		HomeCity:                 input.PatientHomeCity,
		HomeCounty:               input.PatientHomeCounty,
		HomeState:                input.PatientHomeState,
		HomeZip:                  input.PatientHomeZip,
		CirculoConsentFormLink:   input.PatientCirculoConsentFormLink,
		StationMDConsentFormLink: input.PatientStationMDConsentFormLink,
	}

	err := deps.repo.Save(patient)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
