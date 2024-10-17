package idp

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoSonarIdentityProvider struct {
	svc        *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID string
}

func NewCognitoSonarIdentityProvider(userPoolID string) *CognitoSonarIdentityProvider {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return NewCognitoSonarIdentityProviderWithSession(userPoolID, sess)
}

func NewCognitoSonarIdentityProviderWithSession(userPoolID string, sess *session.Session) *CognitoSonarIdentityProvider {
	svc := cognitoidentityprovider.New(sess)

	return &CognitoSonarIdentityProvider{
		svc:        svc,
		userPoolID: userPoolID,
	}
}
