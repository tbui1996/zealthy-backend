package dto

import "github.com/circulohealth/sonar-backend/packages/users/pkg/model"

type ExternalUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Organization *ExternalUserOrganization `json:"organization"`

	Group string `json:"group"`
}

func ExternalUserFromModel(m *model.ExternalUser) *ExternalUser {
	dto := &ExternalUser{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		FirstName: m.FirstName(),
		LastName:  m.LastName(),
		Group:     m.Group(),
	}

	if m.Organization() != nil {
		dto.Organization = ExternalUserOrganizationFromModel(m.Organization())
	}

	return dto
}

func ExternalUsersFromModels(ms []*model.ExternalUser) []*ExternalUser {
	dtos := make([]*ExternalUser, len(ms))
	for i, m := range ms {
		dtos[i] = ExternalUserFromModel(m)
	}

	return dtos
}
