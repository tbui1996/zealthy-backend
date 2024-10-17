package iface

import (
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
)

type ExternalUserOrganizationInsertInput struct {
	Name string
}

type ExternalUserOrganization interface {
	Find(id int) (*model.ExternalUserOrganization, error)
	Insert(record *ExternalUserOrganizationInsertInput) (*model.ExternalUserOrganization, error)
	Update(dm *model.ExternalUserOrganization) (*model.ExternalUserOrganization, error)
	FindAll() ([]*model.ExternalUserOrganization, error)
}
