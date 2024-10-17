package response

import "github.com/circulohealth/sonar-backend/packages/forms/pkg/model"

type SubmittedFormResponse struct {
	Discards    []model.FormDiscard       `json:"discards"`
	Submissions [][]model.InputSubmission `json:"submissions"`
}
