package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	m "github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type PendingChatSessionCreateSuite struct {
	suite.Suite
	req      PendingSessionRequest
	db       *mocks.DynamoDBAPI
	api      *mocks.ApiGatewayManagementApiAPI
	repo     *m.ChatSessionRepository
	sess     *m.ChatSession
	out      map[string]*dynamodb.AttributeValue
	q        *dynamodb.QueryInput
	q1       *dynamodb.QueryInput
	topicReq PendingSessionRequest
}

func (suite *PendingChatSessionCreateSuite) SetupTest() {
	db := new(mocks.DynamoDBAPI)
	api := new(mocks.ApiGatewayManagementApiAPI)
	repo := new(m.ChatSessionRepository)
	sess := new(m.ChatSession)

	suite.db = db
	suite.api = api
	suite.repo = repo
	suite.sess = sess
	suite.q = &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":g": {S: aws.String("internals_general_support")}},
		KeyConditionExpression:    aws.String("CognitoGroup = :g"),
		ProjectionExpression:      aws.String("ConnectionId, UserID"),
		IndexName:                 aws.String("UserGroupIndex"),
		TableName:                 aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	suite.q1 = &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":g": {S: aws.String("internals_program_manager")}},
		KeyConditionExpression:    aws.String("CognitoGroup = :g"),
		ProjectionExpression:      aws.String("ConnectionId, UserID"),
		IndexName:                 aws.String("UserGroupIndex"),
		TableName:                 aws.String(dynamo.SonarInternalWebsocketConnections),
	}

	suite.out = map[string]*dynamodb.AttributeValue{
		"ConnectionId": {
			S: aws.String("1"),
		},
		"UserID": {
			S: aws.String("1"),
		},
	}

	patient := model.Patient{
		Name:        "John",
		LastName:    "Smith",
		Address:     "123 Address",
		InsuranceID: "1234567891234",
		Birthday:    time.Now(),
	}

	desc := model.GENERAL
	suite.req = PendingSessionRequest{
		Repo:   repo,
		Logger: zaptest.NewLogger(suite.T()),
		Request: model.PendingChatSessionCreate{
			UserID:      "1",
			Email:       "test@circulohealth.com",
			Description: &desc,
			Created:     time.Now().Unix(),
			Patient:     &patient,
		},
		DynamoDB: db,
		API:      api,
	}

	topic := "Ian Krieger"
	d := model.CIRCULATOR
	suite.topicReq = PendingSessionRequest{
		Repo:   repo,
		Logger: zaptest.NewLogger(suite.T()),
		Request: model.PendingChatSessionCreate{
			UserID:      "1",
			Email:       "test@circulohealth.com",
			Created:     time.Now().Unix(),
			Topic:       &topic,
			Description: &d,
		},
		DynamoDB: db,
		API:      api,
	}
}

func (suite *PendingChatSessionCreateSuite) TestPendingChatSessionCreate_Success() {
	fmt.Println(suite.q)
	suite.repo.On("CreatePending", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("ID").Return("1")

	suite.db.On("Query", suite.q).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.sess.On("ToResponseDTO").Return(response.ChatSessionResponseDTO{}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, nil)

	actualOut, actualErr := handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "CreatePending", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "ID")
	suite.db.AssertCalled(suite.T(), "Query", suite.q)
	suite.sess.AssertCalled(suite.T(), "ToResponseDTO")
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *PendingChatSessionCreateSuite) TestPendingChatSessionCreate_SuccessNoPatient() {
	suite.repo.On("CreatePending", &suite.topicReq.Request).Return(suite.sess, nil)
	suite.sess.On("ID").Return("1")

	suite.db.On("Query", suite.q1).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.sess.On("ToResponseDTO").Return(response.ChatSessionResponseDTO{}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, nil)

	actualOut, actualErr := handler(suite.topicReq)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "CreatePending", &suite.topicReq.Request)
	suite.sess.AssertCalled(suite.T(), "ID")
	suite.db.AssertCalled(suite.T(), "Query", suite.q1)
	suite.sess.AssertCalled(suite.T(), "ToResponseDTO")
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
}

func (suite *PendingChatSessionCreateSuite) TestPendingChatSessionCreate_FailCreate() {
	suite.repo.On("CreatePending", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	actualOut, actualErr := handler(suite.req)

	suite.Nil(actualOut)
	suite.Error(actualErr)
	suite.repo.AssertCalled(suite.T(), "CreatePending", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "ID")
	suite.db.AssertNotCalled(suite.T(), "Query", suite.q)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	suite.sess.AssertNotCalled(suite.T(), "ToResponseDTO")
}

func (suite *PendingChatSessionCreateSuite) TestPendingChatSessionCreate_FailQuery() {
	suite.repo.On("CreatePending", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("ID").Return("1")

	suite.db.On("Query", suite.q).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	suite.sess.On("ToResponseDTO").Return(response.ChatSessionResponseDTO{}, nil)

	actualOut, actualErr := handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "CreatePending", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "ID")
	suite.db.AssertCalled(suite.T(), "Query", suite.q)
	suite.api.AssertNotCalled(suite.T(), "PostToConnection", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "ToResponseDTO")
}

func (suite *PendingChatSessionCreateSuite) TestPendingChatSessionCreate_FailPost() {
	suite.repo.On("CreatePending", mock.Anything).Return(suite.sess, nil)
	suite.sess.On("ID").Return("1")

	suite.db.On("Query", suite.q).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{suite.out},
	}, nil)

	suite.api.On("PostToConnection", mock.Anything).Return(nil, errors.New("FAKE ERROR, IGNORE"))

	suite.sess.On("ToResponseDTO").Return(response.ChatSessionResponseDTO{}, nil)

	actualOut, actualErr := handler(suite.req)

	suite.NotNil(actualOut)
	suite.NoError(actualErr)
	suite.repo.AssertCalled(suite.T(), "CreatePending", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "ID")
	suite.db.AssertCalled(suite.T(), "Query", suite.q)
	suite.api.AssertCalled(suite.T(), "PostToConnection", mock.Anything)
	suite.sess.AssertCalled(suite.T(), "ToResponseDTO")
}

func TestPendingChatSessionCreate(t *testing.T) {
	suite.Run(t, new(PendingChatSessionCreateSuite))
}
