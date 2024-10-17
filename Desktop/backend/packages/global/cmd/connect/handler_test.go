package main

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func inputHandleRequestBuilder(userID interface{}, connectionID, group string, t *testing.T, dbResponse interface{}, eventPublisherResponse error) HandleRequestInput {
	mockDB := new(dynamo.MockDatabase)
	mockDB.On("Create", mock.Anything).Return(dbResponse)
	mockPublisher := new(mocks.EventPublisher)
	mockPublisher.On("PublishConnectionCreatedEvent", mock.Anything, mock.Anything).Return(eventPublisherResponse)

	return HandleRequestInput{
		DB:             mockDB,
		Event:          eventBuilder(userID, connectionID, group),
		Logger:         zaptest.NewLogger(t),
		EventPublisher: mockPublisher,
	}
}

func TestHandleRequest(t *testing.T) {
	tests := []struct {
		input       HandleRequestInput
		expectedErr *exception.SonarError
	}{
		{
			// valid and successful request
			input:       inputHandleRequestBuilder("userid", "connectionid", "group", t, nil, nil),
			expectedErr: nil,
		},
		{
			// invalid user id
			input:       inputHandleRequestBuilder(nil, "connectionid", "group", t, nil, nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found"),
		},
		{
			// invalid connection id
			input:       inputHandleRequestBuilder("userid", "", "group", t, nil, nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "expected 'connectionID' to be of length > 0"),
		},
		{
			// invalid store to dynamo
			input:       inputHandleRequestBuilder("userid", "connectionid", "group", t, errors.New("storing to dynamo"), nil),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling PutItem for connection item: storing to dynamo"),
		},
		{
			// event failed to publish
			input:       inputHandleRequestBuilder("userid", "connectionid", "group", t, nil, fmt.Errorf("uh-oh")),
			expectedErr: exception.NewSonarError(http.StatusInternalServerError, "uh-oh"),
		},
	}

	for _, test := range tests {
		actualErr := handleRequest(test.input)
		assert.Equal(t, test.expectedErr, actualErr)
	}
}
