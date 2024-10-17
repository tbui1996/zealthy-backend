package iface

import "github.com/circulohealth/sonar-backend/packages/users/pkg/dto"

type InternalUserRepository interface {
	Find(username string) (*dto.InternalUser, error)
}
