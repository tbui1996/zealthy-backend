package validate

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"net/http"
)

type ConnectionContext struct {
	ConnectionID string
	UserID       string
	CognitoGroup string
}

func ConnectionEvent(event events.APIGatewayWebsocketProxyRequest) (*ConnectionContext, *exception.SonarError) {
	userID, userIDExists := event.RequestContext.Authorizer.(map[string]interface{})["userID"].(string)
	if !userIDExists {
		return nil, exception.NewSonarError(http.StatusBadRequest, "expected query params to include 'userID', but was not found")
	}

	connectionID := event.RequestContext.ConnectionID
	if len(connectionID) == 0 {
		return nil, exception.NewSonarError(http.StatusBadRequest, "expected 'connectionID' to be of length > 0")
	}

	group := event.RequestContext.Authorizer.(map[string]interface{})["group"].(string)

	return &ConnectionContext{
		ConnectionID: connectionID,
		UserID:       userID,
		CognitoGroup: group,
	}, nil
}
