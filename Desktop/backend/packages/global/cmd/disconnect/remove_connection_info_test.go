package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func inputRemoveConnectionInfoBuilder(userID, connectionID, group string, dbErr interface{}) RemoveConnectionInfoInput {
	mockDB := new(dynamo.MockDatabase)

	item := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(connectionID),
			},
			"UserID": {
				S: aws.String(userID),
			},
		},
		TableName: aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	mockDB.On("Delete", item).Return(&dynamodb.DeleteItemOutput{}, dbErr)

	return RemoveConnectionInfoInput{
		Context: validate.ConnectionContext{
			ConnectionID: connectionID,
			UserID:       userID,
			CognitoGroup: group,
		},
		DB: mockDB,
	}
}

func TestRemoveConnectionInfo(t *testing.T) {
	tests := []struct {
		input       RemoveConnectionInfoInput
		expectedErr *exception.SonarError
	}{
		{
			// no user id
			input:       inputRemoveConnectionInfoBuilder("", "connectionid", "", errors.New("no UserID")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling DeleteItem for connection item: no UserID"),
		},
		{
			// no connection id
			input:       inputRemoveConnectionInfoBuilder("userid", "", "", errors.New("no ConnectionID")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling DeleteItem for connection item: no ConnectionID"),
		},
		{
			// valid remove, no error
			input:       inputRemoveConnectionInfoBuilder("userid", "connectionid", "group", nil),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		actualErr := removeConnectionInfo(test.input)
		assert.Equal(t, test.expectedErr, actualErr)

		assert.IsType(t, new(dynamo.MockDatabase), test.input.DB)

		mockDB, _ := test.input.DB.(*dynamo.MockDatabase)

		mockDB.AssertCalled(t, "Delete", &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ConnectionId": {
					S: aws.String(test.input.Context.ConnectionID),
				},
				"UserID": {
					S: aws.String(test.input.Context.UserID),
				},
			},
			TableName: aws.String(dynamo.SonarInternalWebsocketConnections),
		})
	}
}
