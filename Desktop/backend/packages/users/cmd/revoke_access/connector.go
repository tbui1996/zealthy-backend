//go:build !test
// +build !test

package main

import (
	"context"
	"encoding/json"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/idp"
	"go.uber.org/zap"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	webSocketUrl := os.Getenv("WEBSOCKET_URL")
	config, err := requestConfig.Must(requestConfig.NewRequestConfig(ctx, event)).ToAPIGatewayV2HTTPRequest()
	if err != nil {
		return exception.NewSonarError(http.StatusInternalServerError, "error setting up config: "+err.Error()).ToAPIGatewayV2HTTPResponse(), nil
	}

	defer logging.SyncLogger(config.Logger)

	config.Logger.Info("revoke_user called")

	var req request.RevokeAccessRequest
	_ = json.Unmarshal([]byte(config.Event.Body), &req)
	err = validateUsername(req.Username)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "invalid request body, body must include 'username'")
	}

	config.Logger = config.Logger.With(zap.String("externalUser", req.Username))
	userPoolID := os.Getenv("USER_POOL_ID")
	api := apigatewaymanagementapi.New(config.Session, &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String(webSocketUrl),
	})
	sonarIDP := idp.NewCognitoSonarIdentityProviderWithSession(userPoolID, config.Session)
	usersDB := dynamo.NewDynamoDatabaseWithSession(dynamo.SonarWebsocketConnections, config.Session)

	reqErr := Handler(&HandlerInput{
		Request:  &req,
		Logger:   config.Logger,
		SonarIDP: sonarIDP,
	})

	if reqErr != nil {
		return reqErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	notificationErr := sendRevokedNotification(&SendRevokedNotificationInput{
		UserID:               req.Username,
		Logger:               config.Logger,
		UsersDB:              usersDB,
		ExternalWebsocketAPI: api,
	})

	if notificationErr != nil {
		return notificationErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("revoke_user complete, sending response")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(connector)
}
