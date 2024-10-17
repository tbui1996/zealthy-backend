package response

import "github.com/circulohealth/sonar-backend/packages/forms/pkg/model"

type FormSentResponse struct {
	Form       model.Form
	Inputs     []model.Input
	FormSentId int
}
