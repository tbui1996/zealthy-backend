package response

import "github.com/circulohealth/sonar-backend/packages/users/pkg/dto"

type Users struct {
	Users           []*dto.ExternalUser `json:"users"`
	PaginationToken *string             `json:"paginationToken"`
}
