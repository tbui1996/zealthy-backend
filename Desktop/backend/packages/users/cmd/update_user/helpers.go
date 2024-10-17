package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"go.uber.org/zap"
)

type UpdateOrganizationDependencies struct {
	User     *model.ExternalUser
	Registry mapper.RegistryAPI
}

type GetSonarWebsocketConnectionDependencies struct {
	Logger *zap.Logger
	Db     dynamodbiface.DynamoDBAPI
}

type SendNotificationDependencies struct {
	Logger *zap.Logger
	Db     dynamodbiface.DynamoDBAPI
	Api    apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

func handleOrganizationUpdate(orgId int, deps UpdateOrganizationDependencies) error {
	if orgId == 0 {
		deps.User.SetOrganization(nil)
		return nil
	}

	organization, err := deps.Registry.ExternalUserOrganization().Find(orgId)
	if err != nil {
		return err
	}

	deps.User.SetOrganization(organization)

	return nil
}

func getSonarWebsocketConnection(email string, deps GetSonarWebsocketConnectionDependencies) (*dynamodb.QueryOutput, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
		KeyConditionExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
		ProjectionExpression: aws.String("ConnectionId"),
	}

	results, err := deps.Db.Query(params)

	deps.Logger.Info(fmt.Sprintf("Query API results: (%s)", results.Items))

	if err != nil {
		return nil, err
	}

	return results, nil
}

func sendConfirmedNotification(email string, deps SendNotificationDependencies) error {
	results, resultErr := getSonarWebsocketConnection(email, GetSonarWebsocketConnectionDependencies{
		Logger: deps.Logger,
		Db:     deps.Db,
	})

	if resultErr != nil {
		return resultErr
	}

	if len(results.Items) == 0 {
		deps.Logger.Debug("no active connections found to send confirmation message")
		return nil
	}

	deps.Logger.Debug("sending confirmation to connection id(s) found")

	var errors []string
	for _, item := range results.Items {
		connectionID := item["ConnectionId"].S
		itemLogger := deps.Logger.With(zap.String("connectionID", *connectionID))
		itemLogger.Info(fmt.Sprintf("Sending message to: (%s)", *connectionID))

		connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: connectionID,
			Data:         []byte("CONFIRMED"),
		}
		_, err := deps.Api.PostToConnection(connectionInput)
		if err != nil {
			errMessage := fmt.Sprintf("sending message to %s: %s", *connectionID, err)
			errors = append(errors, errMessage)
		}
	}

	if len(errors) > 0 {
		// print errors to logs
		deps.Logger.Error(strings.Join(errors, ", "))
		return nil
	}

	return nil
}

func errorResponse(err error, message string, logger *zap.Logger) (events.APIGatewayV2HTTPResponse, error) {
	errMsg := fmt.Errorf("%s %s", message, err.Error())
	logger.Error(errMsg.Error())
	return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg.Error())
}
