package chatHelper

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"go.uber.org/zap"
	"sync"
)

func PostToConnections(connectionItems []model.ConnectionItem, t string, payload interface{}, api apigatewaymanagementapiiface.ApiGatewayManagementApiAPI, logger *zap.Logger) {
	jsonData, err := json.Marshal(&model.InternalUserMessage{
		Type:    t,
		Payload: payload,
	})

	if err != nil {
		logger.Error("Error scanning connection items: %s" + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(connectionItems))
	for _, connectionID := range connectionItems {
		go func(value model.ConnectionItem) {
			defer wg.Done()

			connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(value.ConnectionId),
				Data:         jsonData,
			}

			_, err := api.PostToConnection(connectionInput)

			if err != nil {
				logger.Error(fmt.Sprintf("failed to post pending session to connection (%+v) (%s)", value.ConnectionId, err.Error()))
			}
		}(connectionID)
	}

	wg.Wait()
}
