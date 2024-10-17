package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExternalUser struct {
	ID string
}

func main() {
	lambda.Start(handler)
}

func handler() error {
	logger := logging.Must(logging.NewBasicLogger())
	defer logging.SyncLogger(logger)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	idp := cognitoidentityprovider.New(sess)

	db, err := dao.OpenConnectionWithTablePrefix(dao.Users)
	if err != nil {
		logger.Error("failed to get db connection: " + err.Error())
		return err
	}

	var errs []string
	var paginationToken *string // nil to start and nil once all users have been processed
	baseCase := true
	totalUsersToProcess := 0
	totalUsersProcessed := 0

	for paginationToken != nil || baseCase {
		// do-whiles are hard in golang, paginationToken needs to be nil on first run, then process until it's nil again
		baseCase = false

		usersList, err := getUsers(idp, paginationToken)
		if err != nil {
			logger.Error("failed to get users: " + err.Error())
			errs = append(errs, err.Error())
			break
		}

		totalUsersToProcess += len(usersList.Users)

		dbUsersProcessed := 0
		for _, user := range usersList.Users {
			err := insertUser(db, *user.Username)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to insert user (%s): %s", *user.Username, err.Error()))
				errs = append(errs, err.Error())
			}
			dbUsersProcessed++
		}

		totalUsersProcessed += dbUsersProcessed

		logger.Info(fmt.Sprintf("processed %d of %d users", dbUsersProcessed, len(usersList.Users)))

		paginationToken = usersList.PaginationToken
		if paginationToken == nil {
			logger.Info(fmt.Sprintf("done: processed %d of %d users", totalUsersToProcess, totalUsersToProcess))
		}
	}

	if len(errs) > 0 {
		logger.Error("inserting user ids into db: " + strings.Join(errs, ", "))
		return errors.New("errors inserting user ids into db: " + strings.Join(errs, ", "))
	}

	return nil
}

func insertUser(db *gorm.DB, username string) *exception.SonarError {
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ExternalUser{
		ID: username,
	})

	if result.Error != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "failed to insert user: "+result.Error.Error())
	}

	return nil
}

func getUsers(idp cognitoidentityprovideriface.CognitoIdentityProviderAPI, paginationToken *string) (*cognitoidentityprovider.ListUsersOutput, error) {
	return idp.ListUsers(&cognitoidentityprovider.ListUsersInput{
		UserPoolId:      aws.String(os.Getenv("USER_POOL_ID")),
		PaginationToken: paginationToken,
	})
}
