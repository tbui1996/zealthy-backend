package iface

import "github.com/circulohealth/sonar-backend/packages/users/pkg/model"

type ExternalUser interface {
	Find(id string) (*model.ExternalUser, error)
	FindAll() ([]*model.ExternalUser, error)
	Update(dm *model.ExternalUser) (*model.ExternalUser, error)
}
