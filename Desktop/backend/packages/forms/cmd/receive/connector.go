//go:build !test
// +build !test

package main

import (
	"encoding/json"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"go.uber.org/zap"

	"github.com/circulohealth/sonar-backend/packages/common/iterator"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/response"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"context"
	"fmt"
)

func handleMessage(message events.SQSMessage) error {
	var lf logging.LoggerFields
	logger := logging.Must(lf.FromSQSMessage(message))

	var formTypeResponse response.FormTypeResponse
	if err := json.Unmarshal([]byte(message.Body), &formTypeResponse); err != nil {
		logger.Error("marshalling submit payload: " + err.Error())
		return err
	}

	db, e := dao.OpenConnectionWithTablePrefix(dao.Form)
	if e != nil {
		return e
	}

	switch formTypeResponse.Type {
	case "submit":
		var inputSubmissionRequest request.InputSubmissionRequest
		if err := json.Unmarshal([]byte(formTypeResponse.Message), &inputSubmissionRequest); err != nil {
			logger.Error("marshalling discard payload: " + err.Error())
			return err
		}
		formId := inputSubmissionRequest.FormSentId
		logger = logger.With(zap.Int("formId", formId))
		if len(inputSubmissionRequest.SubmitData) > 0 {
			logger.Debug("saving submit record to DB")
			err := validateRequest(inputSubmissionRequest)
			e = submit(&SubmitInput{
				Db:                     db,
				Err:                    err,
				InputSubmissionRequest: inputSubmissionRequest,
			})
		} else {
			logger.Debug("empty submission request, treating as a discard")
			discardRequest := request.DiscardFormRequest{FormSentId: formId}
			err := validateRequest(discardRequest)
			e = discard(&DiscardInput{
				Db:             db,
				Err:            err,
				DiscardRequest: discardRequest,
				Deleted:        time.Now(),
			})
		}

	case "discard":
		var discardFormRequest request.DiscardFormRequest
		if err := json.Unmarshal([]byte(formTypeResponse.Message), &discardFormRequest); err != nil {
			logger.Error("marshalling discard payload: " + err.Error())
			return err
		}
		err := validateRequest(discardFormRequest)
		e = discard(&DiscardInput{
			Db:             db,
			Err:            err,
			DiscardRequest: discardFormRequest,
			Deleted:        time.Now(),
		})
	default:
		return fmt.Errorf("invalid message type")
	}
	return e
}

func connect(ctx context.Context, event events.SQSEvent) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsIterator := iterator.NewIteratorWithSession(sess, "sonar-service-forms-receive")
	sqsIterator.UseHandler(handleMessage)
	return sqsIterator.Handle(event)

}

func main() {
	lambda.Start(connect)
}
