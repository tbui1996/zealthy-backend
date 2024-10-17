package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func eventBuilder(userID interface{}, connectionID, group string) events.APIGatewayWebsocketProxyRequest {
	return events.APIGatewayWebsocketProxyRequest{
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			Authorizer: map[string]interface{}{
				"userID": userID,
				"group":  group,
			},
			ConnectionID: connectionID,
		},
	}
}

func inputHandleRequestBuilder(userID interface{}, connectionID, group string, t *testing.T, dbErr interface{}) HandleRequestInput {
	mockDB := new(dynamo.MockDatabase)

	_, userIDIsString := userID.(string)
	if !userIDIsString {
		// user id should always be a string, to test when it's not, skip mocking the dynamo call since we can't
		return HandleRequestInput{
			DB:     mockDB,
			Event:  eventBuilder(userID, connectionID, group),
			Logger: zaptest.NewLogger(t),
		}
	}

	item := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(connectionID),
			},
			"UserID": {
				S: aws.String(userID.(string)),
			},
		},
		TableName: aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	mockDB.On("Delete", item).Return(&dynamodb.DeleteItemOutput{}, dbErr)

	return HandleRequestInput{
		DB:     mockDB,
		Event:  eventBuilder(userID, connectionID, group),
		Logger: zaptest.NewLogger(t),
	}
}

func TestHandleRequest(t *testing.T) {
	tests := []struct {
		input       HandleRequestInput
		expectedErr *exception.SonarError
	}{
		{
			// valid and successful request
			input:       inputHandleRequestBuilder("userid", "connectionid", "internals_group", t, nil),
			expectedErr: nil,
		},
		{
			// invalid user id
			input:       inputHandleRequestBuilder(1, "connectionid", "internals_group", t, nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found"),
		},
		{
			// invalid connection id
			input:       inputHandleRequestBuilder("userid", "", "internals_group", t, nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "expected 'connectionID' to be of length > 0"),
		},
		{
			// invalid store to dynamo
			input:       inputHandleRequestBuilder("userid", "connectionid", "internals_group", t, errors.New("storing to dynamo")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling DeleteItem for connection item: storing to dynamo"),
		},
	}

	for _, test := range tests {
		actualErr := handleRequest(test.input)
		assert.Equal(t, test.expectedErr, actualErr)
	}
}
