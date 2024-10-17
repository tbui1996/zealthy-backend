package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"log"
)

type UnconfirmedDisconnectRequest struct {
	Dynamo       dynamodbiface.DynamoDBAPI
	Email        string
	ConnectionId string
}

func Handler(req UnconfirmedDisconnectRequest) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(req.ConnectionId),
			},
			"Email": {
				S: aws.String(req.Email),
			},
		},
		TableName: aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
	}

	_, err := req.Dynamo.DeleteItem(input)
	if err != nil {
		s := fmt.Sprintf("Failed to delete client connection ID %s. Error on DeleteItem: %s", req.ConnectionId, err)
		log.Println(s)
		return errors.New(s)
	}

	return nil
}
