package validate

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/stretchr/testify/assert"
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

func TestValidateAndParseEvent(t *testing.T) {
	tests := []struct {
		inputEvent    events.APIGatewayWebsocketProxyRequest
		expectedEvent *ConnectionContext
		expectedErr   *exception.SonarError
	}{
		{
			// no connection id
			inputEvent:    eventBuilder("userid", "", "group"),
			expectedEvent: nil,
			expectedErr:   exception.NewSonarError(http.StatusBadRequest, "expected 'connectionID' to be of length > 0"),
		},
		{
			// no user id
			inputEvent:    eventBuilder(nil, "connectionid", "group"),
			expectedEvent: nil,
			expectedErr:   exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found"),
		},
		{
			// invalid user id, should be string
			inputEvent:    eventBuilder(1, "connectionid", "group"),
			expectedEvent: nil,
			expectedErr:   exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found"),
		},
		{
			// invalid user id, should be string
			inputEvent:    eventBuilder(1.0, "connectionid", "group"),
			expectedEvent: nil,
			expectedErr:   exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found"),
		},
		{
			// valid user and connection ids
			inputEvent:    eventBuilder("userid", "connectionid", "group"),
			expectedEvent: &ConnectionContext{UserID: "userid", ConnectionID: "connectionid", CognitoGroup: "group"},
			expectedErr:   nil,
		},
	}

	for _, test := range tests {
		actualEvent, actualErr := ConnectionEvent(test.inputEvent)
		assert.Equal(t, test.expectedEvent, actualEvent)
		assert.Equal(t, test.expectedErr, actualErr)
	}
}
