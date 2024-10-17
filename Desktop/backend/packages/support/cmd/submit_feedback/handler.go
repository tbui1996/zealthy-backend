package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"go.uber.org/zap"
)

type SubmitFeedbackRequest struct {
	DynamoDB         dynamodbiface.DynamoDBAPI
	Feedback         request.FeedbackRequest
	Logger           *zap.Logger
	CreatedTimestamp int64
}

type SendFeedbackInput struct {
	FeedbackData []byte
	SesClient    sesiface.SESAPI
	ConfigSet    string
	Template     string
	Domain       string
}

func MarshalFeedbackJSON(name string, feedback request.FeedbackRequest) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":          name,
		"userName":      feedback.UserName,
		"email":         feedback.Email,
		"activity":      feedback.Activity,
		"activityNotes": feedback.ActivityNotes,
		"suggestion":    feedback.Suggestion,
	})
}

func SendFeedback(Input SendFeedbackInput) error {
	sender := "support@" + Input.Domain
	recipient := "madison@circulohealth.com"

	_, err := Input.SesClient.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		ConfigurationSetName: aws.String(Input.ConfigSet),
		Source:               aws.String(sender),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Template:     aws.String(Input.Template),
		TemplateData: aws.String(string(Input.FeedbackData)),
	})

	if err != nil {
		return err
	}

	return nil
}

func Handler(req SubmitFeedbackRequest) error {
	feedback := req.Feedback.Email
	if feedback == nil || *feedback == "" {
		return errors.New("expected an email to be present")
	}

	av, err := dynamodbattribute.MarshalMap(req.Feedback)

	if err != nil {
		errMsg := fmt.Errorf("unable to marshal feedback request to dynamo attributes (%s)", err)
		req.Logger.Error(errMsg.Error())
		return errMsg
	}

	av["createdTimestamp"] = &dynamodb.AttributeValue{
		// nolint gomnd
		N: aws.String(strconv.FormatInt(req.CreatedTimestamp, 10)),
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarFeedback),
	}

	_, err = req.DynamoDB.PutItem(input)

	if err != nil {
		errMsg := fmt.Errorf("unable to save feedback request to dynamo %s", err)
		req.Logger.Error(errMsg.Error())
		return errMsg
	}

	return nil
}
