package response

import (
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
)

type GetFormResponse struct {
	Form   model.Form
	Inputs []model.Input
}
