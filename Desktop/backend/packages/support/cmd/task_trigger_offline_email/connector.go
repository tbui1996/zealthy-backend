//go:build !test
// +build !test

package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/input"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/repo"
)

func connect(context context.Context, userInfo input.UserInfo) error {
	logger, err := logging.NewBasicLogger()
	if err != nil {
		return err
	}
	defer logging.SyncLogger(logger)

	repo := repo.NewOfflineMessageNotificationRepoWithLogger(logger)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sesClient := sesv2.New(sess)

	configSet := os.Getenv("CONFIGURATION_SET_NAME")
	template := os.Getenv("EMAIL_TEMPLATE")
	emailIdentity := os.Getenv("EMAIL_IDENTITY")
	splits := strings.Split(emailIdentity, "/")
	domain := splits[1]

	name, err := MarshalDataJSON(userInfo.FirstName)

	if err != nil {
		return err
	}

	return Handler(userInfo, HandlerDependencies{
		logger,
		repo,
		sesClient,
		configSet,
		template,
		domain,
		name,
	})
}

func main() {
	lambda.Start(connect)
}
