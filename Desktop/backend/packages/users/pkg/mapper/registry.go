package mapper

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// The cost of initializing all mappers instead of just those that are needed
// is likely so low as to not require creating a dynamic registry. Additionally,
// I haven't the experience to figure out how to do that in Go. This is the first
// attempt at a Registry pattern. A second iteration may be generic. That means the registry
// has a generic underlying store rather than specific properties for specific mapper types.

// Initialize as an empty struct, then pass to each mapper, then set each mapper in the registry
type Registry struct {
	externalUser             iface.ExternalUser
	externalUserSQLImpl      externalUserSQLAPI
	externalUserCognitoImpl  externalUserCognitoAPI
	externalUserOrganization iface.ExternalUserOrganization
}

// Logical aggregation of all mapper dependencies so the constructor can construct dependent mappers
type NewRegistryInput struct {
	DB         *gorm.DB
	IDP        cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Logger     *zap.Logger
	UserPoolId *string
}

func NewRegistry(input *NewRegistryInput) *Registry {
	r := &Registry{}

	if input.IDP != nil && input.UserPoolId != nil {
		r.externalUserCognitoImpl = newExternalUserCognito(&newExternalUserCognitoInput{
			idp:        input.IDP,
			userPoolId: *input.UserPoolId,
			logger:     input.Logger,
		})
	}

	if input.DB != nil {
		r.externalUserSQLImpl = newExternalUserSQL(&newExternalUserSQLInput{
			db:     input.DB,
			logger: input.Logger,
		})

		r.externalUserOrganization = newExternalUserOrganization(&newExternalUserOrganizationInput{
			db:     input.DB,
			logger: input.Logger,
		})
	}

	r.externalUser = newExternalUser(&newExternalUserInput{
		registry: r,
		logger:   input.Logger,
	})

	return r
}

func (r *Registry) externalUserCognito() externalUserCognitoAPI {
	return r.externalUserCognitoImpl
}

func (r *Registry) externalUserSQL() externalUserSQLAPI {
	return r.externalUserSQLImpl
}

// Exposes the ExternalUser mapper outside of the package
func (r *Registry) ExternalUser() iface.ExternalUser {
	return r.externalUser
}

// Exposes the ExternalUserOrganization mapper outside of the package
func (r *Registry) ExternalUserOrganization() iface.ExternalUserOrganization {
	return r.externalUserOrganization
}
