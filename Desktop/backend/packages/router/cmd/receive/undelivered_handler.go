package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	common "github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/model"
)

type UndeliveredDelete struct {
	CreatedTimestamp int64
	UserID           string
}

func UndeliveredHandler(db dynamodbiface.DynamoDBAPI, api apigatewaymanagementapiiface.ApiGatewayManagementApiAPI, connectionId string) error {
	out, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.SonarWebsocketConnections),
		Key:       map[string]*dynamodb.AttributeValue{"ConnectionId": {S: aws.String(connectionId)}},
	})

	if err != nil || out.Item == nil {
		return fmt.Errorf("unable to find user id for connection id: %s error: %s", connectionId, err)
	}

	var connectionItem common.ConnectionItem
	err = dynamodbattribute.UnmarshalMap(out.Item, &connectionItem)
	if err != nil {
		return err
	}

	result, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(common.SonarPendingMessages),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {S: aws.String(connectionItem.UserID)},
		},
		ProjectionExpression: aws.String("UserID, CreatedTimestamp, DeleteTimestamp, Message"),
	})

	if err != nil || len(result.Items) == 0 {
		return fmt.Errorf("unable to get pending messages (%s)", err)
	}

	var undelivered []model.UndeliveredMessage
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &undelivered)
	if err != nil {
		return err
	}

	for _, res := range undelivered {
		_, err := api.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(connectionItem.ConnectionId),
			Data:         []byte(res.Message),
		})

		if err != nil {
			continue
		}

		key, _ := dynamodbattribute.MarshalMap(UndeliveredDelete{UserID: res.UserID, CreatedTimestamp: res.CreatedTimestamp})
		_, err = db.DeleteItem(&dynamodb.DeleteItemInput{
			Key:       key,
			TableName: aws.String(common.SonarPendingMessages),
		})

		if err != nil {
			continue
		}
	}

	return nil
}
