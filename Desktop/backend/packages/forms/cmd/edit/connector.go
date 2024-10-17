//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/circulohealth/sonar-backend/packages/common/dao"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/response"
	"github.com/go-playground/validator"

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

	var editFormRequest request.EditForm
	if err := json.Unmarshal([]byte(event.Body), &editFormRequest); err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to parse json body, cannot edit form"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	validate := validator.New()
	err = validate.Struct(editFormRequest)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to validate form request "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to connect to db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	formToEdit, formToEditErr := findFormItem(editFormRequest.ID, db)

	if formToEditErr != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to find form item "+formToEditErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	formToEditErr = editFormItem(&EditFormItemInput{
		Form: formToEdit,
		Req:  editFormRequest,
		Db:   db,
	})

	if formToEditErr != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to update form item "+formToEditErr.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	resp := response.EditFormResponse{Id: strconv.Itoa(formToEdit.ID)}
	body, err := json.Marshal(resp)

	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to marshal json to response"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("edit_form updated form")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
