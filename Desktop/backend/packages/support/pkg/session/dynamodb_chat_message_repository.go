package session

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/uuid"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
)

type DynamoDBChatMessageRepository struct {
	db dynamodbiface.DynamoDBAPI
}

func NewDynamoDBChatMessageRepositoryWithDB(db dynamodbiface.DynamoDBAPI) *DynamoDBChatMessageRepository {
	return &DynamoDBChatMessageRepository{
		db: db,
	}
}

func (repo *DynamoDBChatMessageRepository) GetMessagesForSession(id string) ([]model.ChatMessage, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(dynamo.SonarMessages),
		KeyConditionExpression: aws.String("SessionID = :sessionID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sessionID": {
				S: aws.String(id),
			},
		},
		ProjectionExpression: aws.String("SenderID, Message, CreatedTimestamp, SessionID, ID, FileID"),
	}

	query, err := repo.db.Query(input)
	if err != nil {
		return nil, err
	}

	var messages []model.ChatMessage
	err = dynamodbattribute.UnmarshalListOfMaps(query.Items, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func validateChatRequest(chatRequest request.Chat) error {
	if chatRequest.Message == "" {
		return errors.New("request must include non-empty message")
	}

	if chatRequest.Session == "" {
		return errors.New("request must include session id")
	}

	if chatRequest.Sender == "" {
		return errors.New("request must include a sender id")
	}

	return nil
}

// Stores data into database and assigns a unique ID to it for client uses
func (repo *DynamoDBChatMessageRepository) Create(message request.Chat) (*model.ChatMessage, error) {
	// validate incoming payload
	err := validateChatRequest(message)
	if err != nil {
		return nil, fmt.Errorf("invalid request: (%s)", err)
	}

	dbPayload, err := putItemWithRetry(message, repo.db, 0)

	if err != nil {
		return nil, err
	}

	response := &model.ChatMessage{
		ID:               dbPayload.ID,
		SessionID:        dbPayload.SessionID,
		SenderID:         dbPayload.SenderID,
		Message:          dbPayload.Message,
		CreatedTimestamp: dbPayload.CreatedTimestamp,
	}

	if dbPayload.FileID == "" {
		response.FileID = nil
	} else {
		response.FileID = &dbPayload.FileID
	}

	return response, nil
}

func putItemWithRetry(message request.Chat, db dynamodbiface.DynamoDBAPI, tries int) (*model.Message, error) {
	maxTries := 10
	messageTimestamp := time.Now().Unix()
	uniqueKey := uuid.Create()

	dbPayload := model.Message{
		ID:               uniqueKey,
		SessionID:        message.Session,
		SenderID:         message.Sender,
		Message:          message.Message,
		FileID:           message.File,
		CreatedTimestamp: messageTimestamp,
	}

	av, _ := dynamodbattribute.MarshalMap(dbPayload)

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(dynamo.SonarMessages),
		ConditionExpression: aws.String("attribute_not_exists(SessionID) AND attribute_not_exists(CreatedTimestamp)"),
	}

	_, err := db.PutItem(input)
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException && tries < maxTries {
			// nolint gomnd
			time.Sleep(10 * time.Millisecond)
			// nolint errcheck
			putItemWithRetry(message, db, tries+1)
		} else {
			return nil, err
		}
	}

	return &dbPayload, nil
}
