//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/circulohealth/sonar-backend/packages/common/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, "error setting up config: "+err.Error())
	}

	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("submit_feedback called")

	var feedbackRequest request.FeedbackRequest
	if mErr := json.Unmarshal([]byte(config.Event.Body), &feedbackRequest); mErr != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, mErr.Error())
	}

	db := dynamodb.New(config.Session)
	err = Handler(SubmitFeedbackRequest{
		Feedback:         feedbackRequest,
		DynamoDB:         db,
		Logger:           config.Logger,
		CreatedTimestamp: time.Now().Unix(),
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	svc := ses.New(config.Session)
	configSet := os.Getenv("CONFIGURATION_SET_NAME")
	template := os.Getenv("EMAIL_TEMPLATE")
	emailIdentity := os.Getenv("EMAIL_IDENTITY")
	splits := strings.Split(emailIdentity, "/")

	recipientName := "Madison"
	feedbackData, err := MarshalFeedbackJSON(recipientName, feedbackRequest)

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	err = SendFeedback(SendFeedbackInput{
		FeedbackData: feedbackData,
		SesClient:    svc,
		ConfigSet:    configSet,
		Template:     template,
		Domain:       splits[1],
	})

	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, err.Error())
	}

	config.Logger.Info("submit_feedback complete")

	return response.OKv2()
}

func main() {
	lambda.Start(connector)
}
