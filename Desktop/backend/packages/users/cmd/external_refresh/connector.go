//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}
	defer func() {
		if err := config.Logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	config.Logger.Info("external_refresh called")

	client_id := os.Getenv("CLIENT_ID")

	idp := cognitoidentityprovider.New(config.Session)

	var req request.Refresh
	err = json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "invalid body, expected an object that has a valid 'refreshToken'").ToAPIGatewayV2HTTPResponse(), nil
	}
	if req.RefreshToken == "" {
		return exception.NewSonarError(http.StatusUnprocessableEntity, "invalid body, 'refreshToken' cannot be empty").ToAPIGatewayV2HTTPResponse(), nil
	}

	jsonPayload, reqErr := handler(HandlerInput{
		Idp:          idp,
		ClientID:     client_id,
		RefreshToken: req.RefreshToken,
	})
	if reqErr != nil {
		config.Logger.Error(reqErr.Error())
		return reqErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("external_refresh completed, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonPayload),
	}, nil
}

func main() {
	lambda.Start(connector)
}
