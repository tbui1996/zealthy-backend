package mapper

import "github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"

type RegistryAPI interface {
	ExternalUser() iface.ExternalUser
	ExternalUserOrganization() iface.ExternalUserOrganization
	externalUserSQL() externalUserSQLAPI
	externalUserCognito() externalUserCognitoAPI
}
