//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

	config.Logger.Info("form_delete called")

	if event.PathParameters == nil {
		return exception.NewSonarError(http.StatusBadRequest, "no path params exist").ToAPIGatewayV2HTTPResponse(), nil
	}

	id, ok := event.PathParameters["id"]
	if !ok {
		return exception.NewSonarError(http.StatusBadRequest, "expected path param of {id} to exist").ToAPIGatewayV2HTTPResponse(), nil
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to connect to db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	formToDelete, formToDeleteErr := findFormItem(id, db)

	if formToDeleteErr != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to find form item "+formToDeleteErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	date := time.Now()
	result := addDeleteDate(&AddDeleteDateInput{
		Date:   &date,
		Db:     db,
		FormID: id,
	})

	if result != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to add deleted_at date to form item "+formToDeleteErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	resp := response.DeleteFormResponse{Id: strconv.Itoa(formToDelete.ID)}
	body, err := json.Marshal(resp)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to marshal json to response"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("form_delete deleting form")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
