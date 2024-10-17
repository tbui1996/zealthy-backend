package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AppointmentRepository
	logger *zap.Logger
}

func Handler(input request.CreateAppointmentRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	errMsg := "Patient Id,Appointment Scheduled Time, Appointment Status/Purpose, and Chief Complaint are all required"
	if input.PatientId == "" || input.AppointmentScheduled == "" || input.AppointmentStatus == "" || input.AppointmentPurpose == "" || input.PatientChiefComplaint == "" {
		deps.logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}
	appointment := &model.Appointment{
		PatientId:                     input.PatientId,
		CirculatorDriverFullName:      input.CirculatorDriverFullName,
		AppointmentScheduled:          input.AppointmentScheduled,
		AgencyProviderId:              input.AgencyProviderId,
		AppointmentStatus:             input.AppointmentStatus,
		AppointmentPurpose:            input.AppointmentPurpose,
		AppointmentOtherPurpose:       input.AppointmentOtherPurpose,
		AppointmentNotes:              input.AppointmentNotes,
		PatientDiastolicBloodPressure: input.PatientDiastolicBloodPressure,
		PatientSystolicBloodPressure:  input.PatientSystolicBloodPressure,
		PatientRespirationsPerMinute:  input.PatientRespirationsPerMinute,
		PatientPulseBeatsPerMinute:    input.PatientPulseBeatsPerMinute,
		PatientWeightLbs:              input.PatientWeightLbs,
		PatientChiefComplaint:         input.PatientChiefComplaint,
	}

	err := deps.repo.Save(appointment)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
