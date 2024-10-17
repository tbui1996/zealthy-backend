package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

type PendingSessionRequest struct {
	Repo     iface.ChatSessionRepository
	Logger   *zap.Logger
	Request  model.PendingChatSessionCreate
	DynamoDB dynamodbiface.DynamoDBAPI
	API      apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

func handler(req PendingSessionRequest) ([]byte, error) {
	item := &req.Request
	item.Created = time.Now().Unix()

	desc := model.CIRCULATOR
	if item.Description == nil {
		item.Description = &desc
	} else {
		desc = *item.Description
	}

	chatGroup := model.ChatTypes[desc]

	req.Logger.Debug("creating pending session")
	sess, err := req.Repo.CreatePending(item)
	if err != nil {
		errMessage := fmt.Sprintf("error storing item in dynamo: %+v (%s)", item, err)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	req.Logger = req.Logger.With(zap.String("sessionID", sess.ID()))
	req.Logger.Debug("created pending session, delivering to internal user")

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":g": {S: aws.String(chatGroup)}},
		KeyConditionExpression:    aws.String("CognitoGroup = :g"),
		ProjectionExpression:      aws.String("ConnectionId, UserID"),
		IndexName:                 aws.String("UserGroupIndex"),
		TableName:                 aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	scan, err := req.DynamoDB.Query(input)
	if err != nil {
		req.Logger.Error(fmt.Sprintf("Error scanning connection items: %s", err.Error()))
	}

	if scan != nil {
		var connectionItems []model.ConnectionItem
		req.Logger.Debug(fmt.Sprintf("Found %d items to send to after query", len(scan.Items)))
		err = dynamodbattribute.UnmarshalListOfMaps(scan.Items, &connectionItems)
		if err != nil {
			req.Logger.Error(fmt.Sprintf("unable to unmarshal chat request from message: %+v (%s)", sess, err))
		}
		chatHelper.PostToConnections(connectionItems, "new_pending_session", sess.ToResponseDTO(), req.API, req.Logger)
	}

	bodyBytes, err := json.Marshal(sess.ToResponseDTO())
	if err != nil {
		errMessage := fmt.Errorf("error marshalling session id response: %s", err)
		req.Logger.Error(errMessage.Error())
		return nil, errMessage
	}

	return bodyBytes, nil
}
