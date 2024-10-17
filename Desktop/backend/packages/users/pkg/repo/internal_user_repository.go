package repo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dto"
)

func NewInternalUserRepository() *InternalUserRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	idp := cognitoidentityprovider.New(sess)
	userPoolId := os.Getenv("USER_POOL_ID")

	return &InternalUserRepository{
		idp,
		userPoolId,
	}
}

type InternalUserRepository struct {
	idp        cognitoidentityprovideriface.CognitoIdentityProviderAPI
	userPoolId string
}

func (repo *InternalUserRepository) Find(username string) (*dto.InternalUser, error) {
	record, err := repo.idp.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(repo.userPoolId),
		Username:   aws.String(username),
	})

	if err != nil {
		return nil, err
	}

	result := &dto.InternalUser{
		ID:       *record.Username,
		Username: *record.Username,
	}

	for _, attribute := range record.UserAttributes {
		if attribute.Name != nil {
			switch *attribute.Name {
			case "email":
				result.Email = *attribute.Value
			case "given_name":
				result.FirstName = *attribute.Value
			case "family_name":
				result.LastName = *attribute.Value
			}
		}
	}

	return result, nil
}
