package iface

import (
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
)

type EmailDomainWhitelistRepository interface {
	GetWhitelistDomain(domain string) (*model.EmailDomainWhitelist, *exception.SonarError)
}
