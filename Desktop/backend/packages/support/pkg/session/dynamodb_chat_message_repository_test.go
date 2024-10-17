package session

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DynamoChatMessageRepositoryTestSuite struct {
	suite.Suite
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestGetMessagesForSession_ShouldReturnMessages() {
	db := new(mocks.DynamoDBAPI)

	msg := model.ChatMessage{
		ID:               "1",
		SessionID:        "1",
		SenderID:         "1",
		Message:          "hello",
		CreatedTimestamp: time.Now().Unix(),
		FileID:           nil,
	}

	av, _ := dynamodbattribute.MarshalMap(msg)

	db.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{Items: []map[string]*dynamodb.AttributeValue{av}}, nil)

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	msgs, err := repo.GetMessagesForSession("1")

	suite.NoError(err)
	suite.Equal(1, len(msgs))
	db.AssertCalled(suite.T(), "Query", mock.Anything)
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestGetMessagesForSession_FailQuery() {
	db := new(mocks.DynamoDBAPI)

	db.On("Query", mock.Anything).Return(nil, errors.New("FAKE ERROR"))

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	msgs, err := repo.GetMessagesForSession("1")

	suite.Error(err)
	suite.Nil(msgs)
	db.AssertCalled(suite.T(), "Query", mock.Anything)
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestCreate_ShouldCreate() {
	db := new(mocks.DynamoDBAPI)

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	req := request.Chat{
		Message: "hello",
		Sender:  "1",
		File:    "",
		Session: "1",
	}

	db.On("PutItem", mock.Anything).Return(nil, nil)

	msg, err := repo.Create(req)

	suite.Nil(err)
	suite.Equal("hello", msg.Message)
	db.AssertCalled(suite.T(), "PutItem", mock.Anything)
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestCreate_ShouldCreateWithFile() {
	db := new(mocks.DynamoDBAPI)

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	req := request.Chat{
		Message: "hello",
		Sender:  "1",
		File:    "1",
		Session: "1",
	}

	db.On("PutItem", mock.Anything).Return(nil, nil)

	msg, err := repo.Create(req)

	suite.Nil(err)
	suite.Equal("hello", msg.Message)
	db.AssertCalled(suite.T(), "PutItem", mock.Anything)
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestCreate_ShouldRetryPutItem() {
	db := new(mocks.DynamoDBAPI)

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	req := request.Chat{
		Message: "hello",
		Sender:  "1",
		File:    "1",
		Session: "1",
	}

	db.On("PutItem", mock.Anything).
		Times(1).
		Return(nil, awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "message", errors.New("FAKE ERROR")))
	db.On("PutItem", mock.Anything).Return(nil, nil)

	_, err := repo.Create(req)

	suite.NoError(err)
	db.AssertCalled(suite.T(), "PutItem", mock.Anything)
}

func (suite *DynamoChatMessageRepositoryTestSuite) TestCreate_ShouldFailPutItem() {
	db := new(mocks.DynamoDBAPI)

	repo := NewDynamoDBChatMessageRepositoryWithDB(db)

	req := request.Chat{
		Message: "hello",
		Sender:  "1",
		File:    "1",
		Session: "1",
	}

	db.On("PutItem", mock.Anything).
		Return(nil, awserr.New(dynamodb.ErrCodeDuplicateItemException, "message", errors.New("FAKE ERROR")))

	_, err := repo.Create(req)

	suite.Error(err)
	db.AssertCalled(suite.T(), "PutItem", mock.Anything)
}

func TestDynamoChatMessageRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DynamoChatMessageRepositoryTestSuite))
}
