package forward

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/model"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"go.uber.org/zap"
)

type Forwarder struct {
	// Logger from logger fields sent from message producer
	Logger *zap.Logger
	// Where a message came from
	Source string
	// The external user ids that the message should be delivered to
	Recipients []string
	// Specifies a domain such as forms, chat, or member
	Action string
	// Specifies something that should happen within the "action" (domain)
	Procedure string
	Message   string

	OptOutGuaranteedDelivery bool

	// Forwarder dependencies
	DynamoDB                dynamodbiface.DynamoDBAPI
	ApiGatewayManagementApi apigatewaymanagementapiiface.ApiGatewayManagementApiAPI

	recipientConnections map[string][]model.ConnectionItem

	// the marshalled json payload to be sent through the connection
	data []byte
}

type ForwarderBroadcastDTO struct {
	Sender                  string   `json:"sender"`
	Recipients              []string `json:"recipients"`
	Message                 string   `json:"message"`
	DynamoDB                dynamodbiface.DynamoDBAPI
	ApiGatewayManagementApi apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
	Logger                  *zap.Logger
}

type ForwarderSqsDTO struct {
	Message                 events.SQSMessage
	DynamoDB                dynamodbiface.DynamoDBAPI
	ApiGatewayManagementApi apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
	Logger                  *zap.Logger
}

func NewForwarderFromBroadcast(dto ForwarderBroadcastDTO) (*Forwarder, error) {
	f := Forwarder{
		Logger:                  dto.Logger,
		Source:                  dto.Sender,
		Recipients:              dto.Recipients,
		Message:                 dto.Message,
		Action:                  "broadcast",
		Procedure:               "broadcast",
		DynamoDB:                dto.DynamoDB,
		ApiGatewayManagementApi: dto.ApiGatewayManagementApi,
	}

	return &f, f.Initialize()
}

func NewForwarderFromSQS(dto ForwarderSqsDTO) (*Forwarder, error) {
	actionAttr, ok := dto.Message.MessageAttributes["Action"]
	if !ok {
		return nil, fmt.Errorf("expected action to exist in attributes")
	}

	procedureAttr, ok := dto.Message.MessageAttributes["Procedure"]
	if !ok {
		return nil, fmt.Errorf("expected procedure to exist in attributes")
	}

	sourceAttr, ok := dto.Message.MessageAttributes["Source"]
	if !ok {
		return nil, fmt.Errorf("expected source to exist in attributes")
	}

	f := Forwarder{
		Logger:                  dto.Logger,
		Source:                  *sourceAttr.StringValue,
		Recipients:              []string{},
		Message:                 dto.Message.Body,
		Action:                  *actionAttr.StringValue,
		Procedure:               *procedureAttr.StringValue,
		DynamoDB:                dto.DynamoDB,
		ApiGatewayManagementApi: dto.ApiGatewayManagementApi,
	}

	recipientsAttr, ok := dto.Message.MessageAttributes["Recipients"]
	if ok {
		f.Recipients = strings.Split(*recipientsAttr.StringValue, ",")
		f.Logger = f.Logger.With(zap.String("recipients", *recipientsAttr.StringValue))
	}

	_, ok = dto.Message.MessageAttributes["OptOutGuaranteedDelivery"]
	if ok {
		f.OptOutGuaranteedDelivery = true
		f.Logger = f.Logger.With(zap.Bool("optOutGuaranteedDelivery", f.OptOutGuaranteedDelivery))
	}

	return &f, f.Initialize()
}

func (f *Forwarder) Initialize() error {
	// Marshall data preemptively
	data, err := f.getDataForConnection()
	if err != nil {
		return err
	}

	recipientConnections, err := f.getRecipientConnections()
	if err != nil {
		return err
	}

	f.data = data
	f.recipientConnections = recipientConnections
	return nil
}

func (f *Forwarder) getRecipientConnections() (map[string][]model.ConnectionItem, error) {
	idColumnName := "UserID"
	input, _ := dynamo.GetConnectionItemsInput(dynamo.SonarWebsocketConnections, idColumnName, f.Recipients)

	scan, err := f.DynamoDB.Scan(input)
	if err != nil {
		f.Logger.Error("getting connection items " + err.Error())
		return nil, err
	}

	var connectionItems []model.ConnectionItem
	err = dynamodbattribute.UnmarshalListOfMaps(scan.Items, &connectionItems)
	if err != nil {
		f.Logger.Error("getting marshaling list for connection items " + err.Error())
		return nil, err
	}

	connectionItemsMap := make(map[string][]model.ConnectionItem)

	// make sure every user id is keyed in the map
	for _, userID := range f.Recipients {
		connectionItemsMap[userID] = []model.ConnectionItem{}
	}

	// there can be more than one connection item per user
	for _, cItem := range connectionItems {
		connectionItemsMap[cItem.UserID] = append(connectionItemsMap[cItem.UserID], cItem)
	}

	if len(connectionItemsMap) != len(f.Recipients) {
		f.Logger.Debug("Received a different number of users than recipients. Some aren't connected.")
	}

	return connectionItemsMap, nil
}

func (f *Forwarder) Forward() error {
	f.Logger.Debug("Forward: starting")

	errors := make([]string, 0)
	for user, connectionItems := range f.recipientConnections {
		// There are no connections for the user and the message can't be delivered
		// Store the message in the undelivered dynamodb table
		var err error
		if len(connectionItems) == 0 {
			// If the requester opted out of guaranteed delivery, don't pend messages
			if !f.OptOutGuaranteedDelivery {
				f.Logger.Debug("Forward: pending message from " + f.Source + " for user " + user)
				err = f.pendMessageForUser(user)
			}
		} else {
			f.Logger.Debug("Forward: forwarding message from " + f.Source + " to user " + user)
			err = f.forwardMessageToUser(user)
		}

		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errors, "\n"))
}

func (f *Forwarder) pendMessageForUser(user string) error {
	item := model.UndeliveredMessage{
		UserID:           user,
		CreatedTimestamp: time.Now().Unix(),
		DeleteTimestamp:  time.Now().Add(time.Hour * time.Duration(24)).Unix(), //nolint
		Message:          f.Message,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = f.DynamoDB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarPendingMessages),
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *Forwarder) forwardMessageToUser(user string) error {
	if f.recipientConnections == nil {
		return fmt.Errorf("expected recipientConnections to exist")
	}

	connectionItems, ok := f.recipientConnections[user]

	if !ok {
		return fmt.Errorf("expected user to exist in map")
	}

	errors := make([]string, 0)
	for _, item := range connectionItems {
		err := f.forwardMessageThroughConnection(item.ConnectionId)

		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(errors, "\n"))
}

func (f *Forwarder) getDataForConnection() ([]byte, error) {
	dataUnmarshalled := LoopDataType{
		Action:    f.Action,
		Procedure: f.Procedure,
		Message:   f.Message,
	}

	return json.Marshal(dataUnmarshalled)
}

func (f *Forwarder) forwardMessageThroughConnection(connection string) error {
	connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connection),
		Data:         f.data,
	}

	f.Logger.Debug("Sending message to: " + connection)

	_, err := f.ApiGatewayManagementApi.PostToConnection(connectionInput)
	if err != nil {
		f.Logger.Error("sending message to " + connection + ": " + err.Error())
		return err
	}

	return nil
}
