package appointmenterror

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

type AppointmentRepositoryError struct {
	message string
	code    string
}

func New(message string, code string) *AppointmentRepositoryError {
	return &AppointmentRepositoryError{
		message: message,
		code:    code,
	}
}

func (e *AppointmentRepositoryError) Code() string {
	return e.code
}

func (e *AppointmentRepositoryError) Error() string {
	return e.message
}

func (e *AppointmentRepositoryError) ToSonarError() *exception.SonarError {
	switch e.Code() {
	case NAME_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "name already exists")
	case NOT_FOUND:
		return exception.NewSonarError(http.StatusNotFound, "that appointment was not found")
	case INSURANCE_ID_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "patient with this insurance id already exists")
	case NATIONAL_PROVIDER_ID_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "Agency Provider with this national provider id already exists")
	case DODD_NUMBER_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "Agency Provider with this DoDD number already exists")
	}

	return exception.NewSonarError(http.StatusInternalServerError, "unexpected error occurred while accessing database")
}

func (e *AppointmentRepositoryError) ToSonarApiGatewayResponse() (events.APIGatewayV2HTTPResponse, error) {
	return e.ToSonarError().ToAPIGatewayV2HTTPResponse(), nil
}
