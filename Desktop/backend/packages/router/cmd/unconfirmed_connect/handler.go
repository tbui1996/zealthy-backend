package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
)

type UnconfirmedConnectRequest struct {
	ConnectionId string
	Email        string
	Dynamo       dynamodbiface.DynamoDBAPI
}

func Handler(req UnconfirmedConnectRequest) error {
	av, err := dynamodbattribute.MarshalMap(dynamo.UnconfirmedConnectionItem{
		ConnectionId: req.ConnectionId,
		Email:        req.Email,
	})

	if err != nil {
		return errors.New("error marshaling item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	_, err = req.Dynamo.PutItem(input)

	if err != nil {
		return fmt.Errorf("Failed to connect client %s. Error on PutItem.", req.ConnectionId)
	}

	return nil
}
