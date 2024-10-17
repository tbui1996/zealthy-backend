package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/aws/aws-sdk-go/service/sesv2/sesv2iface"
	"github.com/circulohealth/sonar-backend/packages/common/pretty"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/constants"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/iface"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/input"
	"go.uber.org/zap"
)

type HandlerDependencies struct {
	logger    *zap.Logger
	repo      iface.OfflineMessageNotificationRepo
	sesClient sesv2iface.SESV2API
	configSet string
	template  string
	domain    string
	name      []byte
}

func MarshalDataJSON(name string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name": name,
	})
}

func Handler(userInfo input.UserInfo, deps HandlerDependencies) error {
	deps.logger.Debug(pretty.Sprint(userInfo))

	updated, err := deps.repo.UpdateStatus(userInfo.ID, constants.SENT_NOTIFICATION)

	if err != nil {
		deps.logger.Error(err.Error())
		return err
	}

	if updated {
		deps.logger.Debug(fmt.Sprintf("sending email to: %s", userInfo.Email))
		sender := "support@" + deps.domain

		input := &sesv2.SendEmailInput{
			ConfigurationSetName: aws.String(deps.configSet),
			FromEmailAddress:     aws.String(sender),
			Destination: &sesv2.Destination{
				ToAddresses: []*string{
					aws.String(userInfo.Email),
				},
			},
			Content: &sesv2.EmailContent{
				Template: &sesv2.Template{
					TemplateName: aws.String(deps.template),
					TemplateData: aws.String(string(deps.name)),
				},
			},
			ListManagementOptions: &sesv2.ListManagementOptions{
				ContactListName: aws.String("SonarEmailList"),
				TopicName:       aws.String("newMessageAlert"),
			},
		}

		_, err = deps.sesClient.SendEmail(input)

		if err != nil {
			return err
		}
	}

	return nil
}
