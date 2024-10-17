package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
)

type RemoveConnectionInfoInput struct {
	Context validate.ConnectionContext
	DB      dynamo.Database
}

func removeConnectionInfo(input RemoveConnectionInfoInput) *exception.SonarError {
	item := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(input.Context.ConnectionID),
			},
			"UserID": {
				S: aws.String(input.Context.UserID),
			},
		},
		TableName: aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	_, err := input.DB.Delete(item)
	if err != nil {
		log.Printf("Error deleting dynamo item %s, error: %s", item, err)
		return exception.NewSonarError(http.StatusBadRequest, "calling DeleteItem for connection item: "+err.Error())
	}

	return nil
}
