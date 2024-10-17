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
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator"
)

// HandleRequest form.CreateForm -> *CreateFormResponse
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

	config.Logger.Info("POST /forms: started")

	var createFormRequest request.CreateForm
	if err := json.Unmarshal([]byte(event.Body), &createFormRequest); err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to parse json body, cannot create form"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	validate := validator.New()
	err = validate.Struct(createFormRequest)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to validate form request "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	db, err := dao.OpenConnectionWithTablePrefix(dao.Form)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to connect to db "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	formItem, err := createFormItem(&CreateFormItemInput{
		In:      createFormRequest,
		Db:      db,
		Created: time.Now(),
	})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to create form item "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	err = createInputItems(&CreateInputItemsInput{
		FormId: formItem.ID,
		Form:   createFormRequest,
		Db:     db})
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not create inputs for form."+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	resp := response.CreateFormResponse{Id: strconv.Itoa(formItem.ID)}
	body, err := json.Marshal(resp)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "unable to marshal json to response"+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("POST /forms: form created")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(connect)
}
