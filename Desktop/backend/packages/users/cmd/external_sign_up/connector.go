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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dao"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
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

	idp := cognitoidentityprovider.New(config.Session)
	db := dynamodb.New(config.Session)
	repo := dao.NewDynamoDBEmailDomainWhitelistRepository(db)
	clientID := os.Getenv("CLIENT_ID")

	var req request.ExternalSignUpRequest
	err = json.Unmarshal([]byte(config.Event.Body), &req)
	if err != nil {
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, "invalid request body "+err.Error())
	}

	config.Logger.Info("external_sign_up called")

	fullNameSplit := strings.Fields(req.FullName)

	output, sErr := handler(HandlerInput{
		Event:            event,
		Idp:              idp,
		Repo:             repo,
		ClientID:         clientID,
		Logger:           config.Logger,
		FamilyName:       fullNameSplit[len(fullNameSplit)-1],
		GivenName:        fullNameSplit[0],
		OrganizationName: req.OrganizationName,
	})

	if sErr != nil {
		return sErr.ToAPIGatewayV2HTTPResponse(), nil
	}

	config.Logger.Info("external_sign_up complete, sending response", zap.String("userID", *output.UserSub))

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	lambda.Start(connector)
}
