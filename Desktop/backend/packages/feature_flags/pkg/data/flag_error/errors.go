package flagerror

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
)

type FeatureFlagRepositoryError struct {
	message string
	code    string
}

func New(message string, code string) *FeatureFlagRepositoryError {
	return &FeatureFlagRepositoryError{
		message: message,
		code:    code,
	}
}

func (e *FeatureFlagRepositoryError) Code() string {
	return e.code
}

func (e *FeatureFlagRepositoryError) Error() string {
	return e.message
}

func (e *FeatureFlagRepositoryError) ToSonarError() *exception.SonarError {
	switch e.Code() {
	case KEY_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "flagKey already exists")
	case NAME_CONFLICT:
		return exception.NewSonarError(http.StatusConflict, "name already exists")
	case NOT_FOUND:
		return exception.NewSonarError(http.StatusNotFound, "that flag was not found")
	}
	return exception.NewSonarError(http.StatusInternalServerError, "unexpected error occurred while accessing database")
}

func (e *FeatureFlagRepositoryError) ToSonarApiGatewayResponse() (events.APIGatewayV2HTTPResponse, error) {
	return e.ToSonarError().ToAPIGatewayV2HTTPResponse(), nil
}
