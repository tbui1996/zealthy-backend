//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
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

	config.Logger.Info("Get /forms/{id}/response: started")
	if event.PathParameters == nil {
		return exception.NewSonarError(http.StatusBadRequest, "expected path parameters to exist").ToAPIGatewayV2HTTPResponse(), nil
	}
	id, ok := event.PathParameters["id"]
	if !ok {
		return exception.NewSonarError(http.StatusBadRequest, "parameter path {id} was not found").ToAPIGatewayV2HTTPResponse(), nil
	}
	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not connect to database, unable to retrieve form "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	sents, err := findFormSent(id, db)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not retrieve sent form "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	formSents := make([]int, len(sents))
	for i, v := range sents {
		formSents[i] = v.ID
	}

	submit, err := findSubmitByFormSent(formSents, db)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not get form submissions "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	discards, inputs, err := getDiscardAndSubmitValues(&DiscardAndSubmitInput{
		Db:        db,
		FormSents: formSents,
		Submit:    submit,
	})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not discard and submit "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	body, err := json.Marshal(response.SubmittedFormResponse{Discards: discards, Submissions: inputs})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to parse JSON response "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	config.Logger.Info("Get /forms/{id}/response: complete")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
