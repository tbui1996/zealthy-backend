package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"go.uber.org/zap"
)

type LoopSendRequest struct {
	Logger   *zap.Logger
	Repo     iface.ChatSessionRepository
	Message  request.Chat
	DynamoDB dynamodbiface.DynamoDBAPI
	API      apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

type LoopSendResponse struct {
	ID                 string  `json:"id"`
	SessionID          string  `json:"sessionID"`
	SenderID           string  `json:"senderID"`
	Message            string  `json:"message"`
	CreatedTimestamp   int64   `json:"createdTimestamp"`
	FileID             *string `json:"fileID"`
	IsSupportAvailable bool    `json:"isSupportAvailable"`
}

func errorResponse(message string, err error, logger *zap.Logger) error {
	errMsg := fmt.Errorf("%s\n%s", message, err.Error())
	logger.Error(errMsg.Error())
	return errMsg
}

func loopResponse(messageResponse *model.ChatMessage, supportAvailable bool) ([]byte, error) {
	msgResp := *messageResponse
	return json.Marshal(LoopSendResponse{
		ID:                 msgResp.ID,
		SessionID:          msgResp.SessionID,
		SenderID:           msgResp.SenderID,
		Message:            msgResp.Message,
		CreatedTimestamp:   msgResp.CreatedTimestamp,
		FileID:             msgResp.FileID,
		IsSupportAvailable: supportAvailable,
	})
}

func Handler(req LoopSendRequest) ([]byte, error) {
	sess, err := req.Repo.GetEntityWithUsers(req.Message.Session)

	if err != nil {
		return nil, errorResponse(fmt.Sprintf("error getting session %s", req.Message.Session), err, req.Logger)
	}

	req.Logger.Debug("received session, appending message to session")
	messageResponse, err := sess.AppendRequestMessage(req.Message)
	if err != nil {
		return nil, errorResponse("unable to store message in DB", err, req.Logger)
	}

	req.Logger = req.Logger.With(zap.String("id", messageResponse.ID))
	req.Logger.Debug("appended message to session, delivering to internal user")

	var input = &dynamodb.QueryInput{
		ProjectionExpression: aws.String("ConnectionId, UserID"),
		TableName:            aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	if sess.IsPending() {
		chatType := sess.Type()
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{":g": {S: aws.String(model.ChatTypes[chatType])}}
		input.KeyConditionExpression = aws.String("CognitoGroup = :g")
		input.IndexName = aws.String("UserGroupIndex")
	} else {
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{":user": {S: aws.String(sess.InternalUserID())}}
		input.KeyConditionExpression = aws.String("UserID = :user")
	}

	scan, err := req.DynamoDB.Query(input)
	if err != nil {
		req.Logger.Error(fmt.Sprintf("error scanning connection items: %s", err))
	}

	var connectionItems []model.ConnectionItem
	if scan != nil {
		req.Logger.Debug(fmt.Sprintf("found %d items to send to after query", len(scan.Items)))
		err = dynamodbattribute.UnmarshalListOfMaps(scan.Items, &connectionItems)
		if err != nil {
			req.Logger.Error(fmt.Sprintf("error unmarshaling connection items input: %s", err))
		}
		chatHelper.PostToConnections(connectionItems, "message", messageResponse, req.API, req.Logger)
	}

	response, err := loopResponse(messageResponse, len(connectionItems) > 0)
	if err != nil {
		return nil, errorResponse("failed to marshall message response", err, req.Logger)
	}
	return response, nil
}
