package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
)

func main() {
	lambda.Start(handler)
}

func handler() error {
	logger := logging.Must(logging.NewBasicLogger())
	defer logging.SyncLogger(logger)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sesClient := sesv2.New(sess)

	_, err := sesClient.CreateContactList(&sesv2.CreateContactListInput{
		ContactListName: aws.String("SonarEmailList"),
		Topics: []*sesv2.Topic{
			{
				DefaultSubscriptionStatus: aws.String("OPT_IN"),
				Description:               aws.String("Never miss a Circulo Chat message! Sonar will send you an email if you receive messages while offline."),
				DisplayName:               aws.String("New Message Alert"),
				TopicName:                 aws.String("newMessageAlert"),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
