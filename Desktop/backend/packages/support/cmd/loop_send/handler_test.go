package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	m "github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

type LoopSendSuite struct {
	suite.Suite
	req  LoopSendRequest
	db   *mocks.DynamoDBAPI
	api  *mocks.ApiGatewayManagementApiAPI
	repo *m.ChatSessionRepository
	sess *m.ChatSession
	out  map[string]*dynamodb.AttributeValue
	msg  *model.ChatMessage
	t    model.ChatType
}

func (suite *LoopSendSuite) SetupTest() {
	db := new(mocks.DynamoDBAPI)
	api := new(mocks.ApiGatewayManagementApiAPI)
	repo := new(m.ChatSessionRepository)
	sess := new(m.ChatSession)

	suite.db = db
	suite.api = api
	suite.repo = repo
	suite.sess = sess
	suite.t = model.GENERAL

	suite.out = map[string]*dynamodb.AttributeValue{
		"ConnectionId": {
			S: aws.String("1"),
		},
		"UserID": {
			S: aws.String("1"),
		},
	}

	suite.msg = &model.ChatMessage{
		ID:               "1",
		SessionID:        "1",
		SenderID:         "1",
		Message:          "Hello",
		CreatedTimestamp: time.Now().Unix(),
		FileID:           nil,
	}

	suite.req = LoopSendRequest{
		Repo:   repo,
		Logger: zaptest.NewLogger(suite.T()),
		Message: request.Chat{
			Session: suite.msg.SessionID,
			Sender:  suite.msg.SenderID,
			Message: suite.msg.Message,
			File:    "",
		},
		DynamoDB: db,
		API:      api,
	}
}

func (suite *LoopSendSuite) TestLoopSend_Success() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("Type").Return(suite.t, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, nil)
	suite.sess.On("IsPending", mock.Anything).Return(true)

	suite.db.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, nil)

	actualOut, actualErr := Handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "IsPending")
	suite.db.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_SuccessNotPending() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, nil)
	suite.sess.On("IsPending", mock.Anything).Return(false)
	suite.sess.On("InternalUserID", mock.Anything).Return("1")

	suite.db.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, nil)

	actualOut, actualErr := Handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "IsPending")
	suite.sess.AssertCalled(suite.T(), "InternalUserID")
	suite.db.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_SuccessWithNoConnectionItems() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, nil)
	suite.sess.On("IsPending", mock.Anything).Return(false)
	suite.sess.On("InternalUserID", mock.Anything).Return("1")

	suite.db.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}, nil)

	actualOut, actualErr := Handler(suite.req)

	var resp LoopSendResponse
	_ = json.Unmarshal(actualOut, &resp)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.False(resp.IsSupportAvailable)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "IsPending")
	suite.sess.AssertCalled(suite.T(), "InternalUserID")
	suite.db.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_FailGetEntity() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	actualOut, actualErr := Handler(suite.req)

	suite.Nil(actualOut)
	suite.Error(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "IsPending")
	suite.sess.AssertNotCalled(suite.T(), "InternalUserID")
	suite.db.AssertNotCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_FailAppendRequest() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, errors.New("FAKE ERROR, IGNORE"))

	actualOut, actualErr := Handler(suite.req)

	suite.Nil(actualOut)
	suite.Error(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "IsPending")
	suite.sess.AssertNotCalled(suite.T(), "InternalUserID")
	suite.db.AssertNotCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_FailOnQuery() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, nil)
	suite.sess.On("IsPending", mock.Anything).Return(false)
	suite.sess.On("InternalUserID", mock.Anything).Return("1")

	suite.db.On("Query", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	actualOut, actualErr := Handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "IsPending")
	suite.sess.AssertCalled(suite.T(), "InternalUserID")
	suite.db.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *LoopSendSuite) TestLoopSend_FailOnPost() {
	suite.repo.On("GetEntityWithUsers", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("AppendRequestMessage", mock.Anything).Return(suite.msg, nil)
	suite.sess.On("IsPending", mock.Anything).Return(false)
	suite.sess.On("InternalUserID", mock.Anything).Return("1")

	suite.db.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	actualOut, actualErr := Handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "GetEntityWithUsers", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "Type")
	suite.sess.AssertCalled(suite.T(), "AppendRequestMessage", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "IsPending")
	suite.sess.AssertCalled(suite.T(), "InternalUserID")
	suite.db.AssertCalled(suite.T(), "Query", mock.Anything)
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func TestLoopSend(t *testing.T) {
	suite.Run(t, new(LoopSendSuite))
}
