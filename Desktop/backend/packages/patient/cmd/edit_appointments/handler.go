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
	repo   iface.AppointmentRepository
	logger *zap.Logger
}

func Handler(input request.EditAppointmentRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	appointment, err := deps.repo.Find(input.AppointmentId)
	if err != nil {
		deps.logger.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
	}
	now := time.Now()
	appointment.AppointmentNotes = input.AppointmentNotes
	appointment.LastModifiedTimestamp = &now
	appointment.AppointmentPurpose = input.AppointmentPurpose
	appointment.AppointmentOtherPurpose = input.AppointmentOtherPurpose
	appointment.PatientChiefComplaint = input.PatientChiefComplaint
	appointment.AgencyProviderId = input.AgencyProviderId
	appointment.CirculatorDriverFullName = input.CirculatorDriverFullName
	appointment.PatientDiastolicBloodPressure = input.PatientDiastolicBloodPressure
	appointment.PatientSystolicBloodPressure = input.PatientSystolicBloodPressure
	appointment.PatientRespirationsPerMinute = input.PatientRespirationsPerMinute
	appointment.PatientPulseBeatsPerMinute = input.PatientPulseBeatsPerMinute
	appointment.PatientWeightLbs = input.PatientWeightLbs
	appointment.AppointmentScheduled = input.AppointmentScheduled
	appointment.PatientId = input.PatientId

	if input.AppointmentStatus != appointment.AppointmentStatus {
		appointment.AppointmentStatusChangedOn = &now
		appointment.AppointmentStatus = input.AppointmentStatus
	}

	err = deps.repo.Save(appointment)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
