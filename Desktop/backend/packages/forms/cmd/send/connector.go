//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/repo"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func connect(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()
	config.Logger.Info("edit_form called")
	if event.PathParameters == nil {
		return exception.NewSonarError(http.StatusBadRequest, "no path params exist").ToAPIGatewayV2HTTPResponse(), nil
	}
	id, ok := event.PathParameters["id"]
	if !ok {
		return exception.NewSonarError(http.StatusBadRequest, "expected path param of {id} to exist").ToAPIGatewayV2HTTPResponse(), nil
	}
	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to open database connection"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	repository := repo.NewRepository(db)
	inputs, err := repository.Inputs(id)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not retrieve form inputs from db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	form, err := repository.Form(id)
	if err != nil || form == nil || form.DateClosed != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to get form from db").ToAPIGatewayV2HTTPResponse(), nil
	}

	sent, err := createSendRecord(&CreateSendRecordInput{
		FormID: form.ID,
		Sent:   time.Now(),
		Db:     db,
	})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to save form sent record").ToAPIGatewayV2HTTPResponse(), nil
	}

	body, err := json.Marshal(&response.FormSentResponse{Form: *form, Inputs: inputs, FormSentId: sent.ID})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to marshal JSON response "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	client := router.NewClientWithConfigWithSession(&router.Config{
		SendQueueName:    "sonar-service-forms-send",
		ReceiveQueueName: "sonar-service-forms-receive",
	}, config.Session)
	err = sendThroughRouter(client, string(body))
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to send message via sqs queue "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}

func main() {
	lambda.Start(connect)
}
