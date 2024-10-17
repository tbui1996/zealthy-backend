package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/connection"
	connectionDto "github.com/circulohealth/sonar-backend/packages/common/connection/dto"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"go.uber.org/zap"
)

func connect(context context.Context, userID string) (connectionDto.UserOnlineStatus, error) {
	logger, err := logging.NewBasicLogger()
	if err != nil {
		return connectionDto.UserOnlineStatus{}, err
	}
	defer logging.SyncLogger(logger)

	logger.With(zap.String("userID", userID))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	onlineChecker := connection.OnlineChecker{
		TableName: dynamo.SonarInternalWebsocketConnections,
		DB:        db,
		Logger:    logger,
	}

	return onlineChecker.IsUserOnline(userID)
}

func main() {
	lambda.Start(connect)
}
