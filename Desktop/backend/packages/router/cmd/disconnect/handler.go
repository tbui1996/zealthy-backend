package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
)

type DisconnectRequest struct {
	UserId       string
	ConnectionId string
	Dynamo       dynamodbiface.DynamoDBAPI
}

func Handler(req DisconnectRequest) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(req.ConnectionId),
			},
			"UserID": {
				S: aws.String(req.UserId),
			},
		},
		TableName: aws.String(dynamo.SonarWebsocketConnections),
	}

	_, err := req.Dynamo.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("Failed to delete client connection ID %s. Error on DeleteItem: %s", req.ConnectionId, err)
	}

	return nil
}
