package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/events/eventconstants"
	"github.com/circulohealth/sonar-backend/packages/common/events/iface"
)

type ConnectRequest struct {
	UserID         string
	ConnectionId   string
	Dynamo         dynamodbiface.DynamoDBAPI
	EventPublisher iface.EventPublisher
}

func Handler(req ConnectRequest) error {
	av, err := dynamodbattribute.MarshalMap(dynamo.ConnectionItem{
		ConnectionId: req.ConnectionId,
		UserID:       req.UserID,
	})

	if err != nil {
		return errors.New("error marshaling connection item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarWebsocketConnections),
	}

	_, err = req.Dynamo.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to connect client %s. Error on PutItem", req.ConnectionId)
	}

	err = req.EventPublisher.PublishConnectionCreatedEvent(req.UserID, eventconstants.ROUTER_SERVICE)

	return err
}
