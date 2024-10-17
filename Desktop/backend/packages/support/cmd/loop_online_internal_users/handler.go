package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	db     dynamodbiface.DynamoDBAPI
	logger *zap.Logger
}

func Handler(deps HandlerDeps) (*[]response.OnlineUserResponse, error) {
	dynamoInput := &dynamodb.ScanInput{
		TableName: aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	result, err := deps.db.Scan(dynamoInput)

	if err != nil {
		deps.logger.Error(fmt.Sprintf("Unable to get items from dynamo: %s", err.Error()))
		return nil, err
	}

	var connectionItems []model.ConnectionItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &connectionItems)
	if err != nil {
		deps.logger.Error(fmt.Sprintf("getting marshaling list for connection items: %s", err.Error()))
		return nil, err
	}

	alreadySeenUsers := make(map[string]model.ConnectionItem)

	results := make([]response.OnlineUserResponse, 0)

	for _, item := range connectionItems {
		if _, ok := alreadySeenUsers[item.UserID]; ok {
			continue
		}
		alreadySeenUsers[item.UserID] = item

		results = append(results, response.OnlineUserResponse{
			UserID: item.UserID,
		})
	}

	return &results, nil
}
