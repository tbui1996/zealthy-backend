package mapper

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	userPoolId := "1"
	input := NewRegistryInput{
		Logger:     zaptest.NewLogger(t),
		UserPoolId: &userPoolId,
		DB:         theDB,
		IDP:        new(mocks.CognitoIdentityProviderAPI),
	}

	r := NewRegistry(&input)
	assert.NotNil(t, r.externalUserCognito())
	assert.NotNil(t, r.externalUserSQL())
	assert.NotNil(t, r.ExternalUser())
	assert.NotNil(t, r.ExternalUserOrganization())
}
