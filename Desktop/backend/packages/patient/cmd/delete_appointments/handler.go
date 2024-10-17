package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AppointmentRepository
	logger *zap.Logger
}

func Handler(input request.DeleteAppointmentRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	if input.AppointmentId == "" {
		errMsg := "appointment_id is required"
		deps.logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}

	appointment, err := deps.repo.Find(input.AppointmentId)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	newError := deps.repo.Delete(appointment)

	if newError != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
