package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/response"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AppointmentRepository
	logger *zap.Logger
}

func Handler(deps HandlerDeps) (*[]response.AppointmentResponse, error) {
	appointments, err := deps.repo.FindAll()

	if err != nil {
		deps.logger.Error(err.Error())
		return nil, fmt.Errorf("unable to find records in database")
	}

	responses := make([]response.AppointmentResponse, len(*appointments))

	for i := range *appointments {
		appointment := (*appointments)[i]
		responses[i] = response.AppointmentResponse{
			AppointmentId:                 appointment.AppointmentId,
			PatientId:                     appointment.PatientId,
			AgencyProviderId:              appointment.AgencyProviderId,
			CirculatorDriverFullName:      appointment.CirculatorDriverFullName,
			AppointmentCreated:            appointment.AppointmentCreated,
			AppointmentScheduled:          appointment.AppointmentScheduled,
			AppointmentStatus:             appointment.AppointmentStatus,
			AppointmentStatusChangedOn:    appointment.AppointmentStatusChangedOn,
			AppointmentPurpose:            appointment.AppointmentPurpose,
			AppointmentOtherPurpose:       appointment.AppointmentOtherPurpose,
			AppointmentNotes:              appointment.AppointmentNotes,
			PatientDiastolicBloodPressure: appointment.PatientDiastolicBloodPressure,
			PatientSystolicBloodPressure:  appointment.PatientSystolicBloodPressure,
			PatientRespirationsPerMinute:  appointment.PatientRespirationsPerMinute,
			PatientPulseBeatsPerMinute:    appointment.PatientPulseBeatsPerMinute,
			PatientWeightLbs:              appointment.PatientWeightLbs,
			PatientChiefComplaint:         appointment.PatientChiefComplaint,
			CreatedTimestamp:              appointment.CreatedTimestamp,
			LastModifiedTimestamp:         appointment.LastModifiedTimestamp,
			FirstName:                     appointment.FirstName,
			MiddleName:                    appointment.MiddleName,
			LastName:                      appointment.LastName,
			ProviderFullName:              appointment.ProviderFullName,
			Suffix:                        appointment.Suffix,
			DateOfBirth:                   &appointment.DateOfBirth,
			PrimaryLanguage:               appointment.PrimaryLanguage,
			PreferredGender:               appointment.PreferredGender,
			EmailAddress:                  appointment.EmailAddress,
			HomeAddress1:                  appointment.HomeAddress1,
			HomeAddress2:                  appointment.HomeAddress2,
			HomeCity:                      appointment.HomeCity,
			HomeState:                     appointment.HomeState,
			HomeZip:                       appointment.HomeZip,
			SignedCirculoConsentForm:      appointment.SignedCirculoConsentForm,
			CirculoConsentFormLink:        appointment.CirculoConsentFormLink,
			SignedStationMDConsentForm:    appointment.SignedStationMDConsentForm,
			StationMDConsentFormLink:      appointment.StationMDConsentFormLink,
			CompletedGoSheet:              appointment.CompletedGoSheet,
			MarkedAsActive:                appointment.MarkedAsActive,
			NationalProviderId:            appointment.NationalProviderId,
			BusinessName:                  appointment.BusinessName,
			BusinessTIN:                   appointment.BusinessTIN,
			BusinessAddress1:              appointment.BusinessAddress1,
			BusinessAddress2:              appointment.BusinessAddress2,
			BusinessCity:                  appointment.BusinessCity,
			BusinessState:                 appointment.BusinessState,
			BusinessZip:                   appointment.BusinessZip,
			HomePhone:                     appointment.HomePhone,
			HomeLivingArrangement:         appointment.HomeLivingArrangement,
			HomeCounty:                    appointment.HomeCounty,
			InsuranceId:                   appointment.InsuranceId,
		}
	}

	return &responses, nil
}
