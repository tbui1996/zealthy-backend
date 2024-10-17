package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dto"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
)

func Connect(context context.Context, userID string) (dto.ExternalUser, error) {
	logger, err := logging.NewBasicLogger()
	if err != nil {
		return dto.ExternalUser{}, err
	}

	defer logging.SyncLogger(logger)

	logger.Debug(fmt.Sprintf("getting user info for: %s", userID))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)
	if err != nil {
		return dto.ExternalUser{}, err
	}
	idp := cognitoidentityprovider.New(sess)
	userPoolID := os.Getenv("USER_POOL_ID")

	registry := mapper.NewRegistry(&mapper.NewRegistryInput{
		DB:         db,
		IDP:        idp,
		UserPoolId: &userPoolID,
		Logger:     logger,
	})

	logger.Debug("getting mapper")
	mapper := registry.ExternalUser()

	logger.Debug(fmt.Sprintf("finding user info for: %s", userID))
	user, err := mapper.Find(userID)
	logger.Debug(fmt.Sprintf("found user info for: %s", userID))

	if err != nil {
		logger.Error(err.Error())
		return dto.ExternalUser{}, err
	}

	logger.Debug("Parsing DTO")
	dto := dto.ExternalUserFromModel(user)

	logger.Debug("Done")
	return *dto, nil
}

func main() {
	lambda.Start(Connect)
}
