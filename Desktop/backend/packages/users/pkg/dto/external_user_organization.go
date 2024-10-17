package dto

import "github.com/circulohealth/sonar-backend/packages/users/pkg/model"

type ExternalUserOrganization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ExternalUserOrganizationFromModel(m *model.ExternalUserOrganization) *ExternalUserOrganization {
	return &ExternalUserOrganization{
		ID:   m.ID,
		Name: m.Name(),
	}
}
