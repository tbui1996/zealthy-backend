package request

import (
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
)

type CreateForm struct {
	// Title must be at least 4 characters
	Title string `json:"title" validate:"required,min=1"`
	// description is not required
	Description string `json:"description"`
	// creator is required
	Creator string `json:"creator" validate:"required,min=1"`
	// creator is required
	CreatorId string `json:"creatorId" validate:"required,min=1"`
	// Must have at least 1 input
	Inputs []model.Input `json:"inputs" validate:"required,min=1"`
}
