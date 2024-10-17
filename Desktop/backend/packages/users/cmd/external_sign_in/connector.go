//go:build !test
// +build !test

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/externalAuth"
	"go.uber.org/zap"
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

	config.Logger.Info("external_sign_in called")

	clientID := os.Getenv("CLIENT_ID")
	userPoolID := os.Getenv("USER_POOL_ID")
	idp := cognitoidentityprovider.New(config.Session)
	email, err := externalAuth.GetEmailFromToken(config.Event.Headers)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger = config.Logger.With(zap.String("email", email))

	body, sErr := handler(HandlerInput{
		Logger:   config.Logger,
		Idp:      idp,
		Email:    email,
		ClientID: clientID,
		PoolID:   userPoolID,
	})
	if sErr != nil {
		config.Logger.Error("failed to sign in: " + sErr.Error())
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("external_sign_in complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(connector)
}
