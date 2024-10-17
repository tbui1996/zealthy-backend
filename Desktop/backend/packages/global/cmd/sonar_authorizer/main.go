//go:build !test
// +build !test

package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

func HandleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	groupNamePrefix := os.Getenv("GROUP_NAME_PREFIX")
	jwksUrl := os.Getenv("JWKS_URL")

	logger := logging.Must(logging.NewLoggerFromEvent(event))
	defer func() {
		if err := logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
			log.Println(err)
		}
	}()

	logger = logger.With(zap.String("groupNamePrefix", groupNamePrefix))
	logger.Info("sonar_authorizer started")

	db := &dynamo.DynamoDatabase{
		TableName: dynamo.SonarGroupPolicy,
	}

	keySet, err := jwk.Fetch(ctx, jwksUrl)
	if err != nil {
		errMessage := "failed to get public JSON web key set to validate tokens: " + err.Error()
		logger.Error(errMessage)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New(errMessage)
	}

	input := HandleRequestInput{
		Context:         ctx,
		Event:           &event,
		GroupNamePrefix: groupNamePrefix,
		Logger:          logger,
		DB:              db,
		JwkSet:          keySet,
		Jwt:             &authorizer.Jwt{},
	}

	response, err := handleRequest(input)
	if err != nil {
		logger.Error("building response for user: " + err.Error())
	}

	return response, err
}

func main() {
	lambda.Start(HandleRequest)
}
