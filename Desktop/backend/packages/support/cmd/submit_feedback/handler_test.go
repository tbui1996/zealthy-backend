package main

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/suite"
)

type SubmitFeedbackTestSuite struct {
	suite.Suite
	Request    SubmitFeedbackRequest
	PutItem    *dynamodb.PutItemInput
	DynamoMock *mocks.DynamoDBAPI
}

func (suite *SubmitFeedbackTestSuite) SetupTest() {
	db := new(mocks.DynamoDBAPI)
	suite.DynamoMock = db

	email := "test@circulohealth.com"
	feedback := request.FeedbackRequest{
		Email:         &email,
		Activity:      "chatting",
		ActivityNotes: "wanted to chat but couldn't",
		Suggestion:    "make chat work",
	}

	u := time.Now().Unix()
	suite.Request = SubmitFeedbackRequest{
		DynamoDB:         db,
		Logger:           zaptest.NewLogger(suite.T()),
		Feedback:         feedback,
		CreatedTimestamp: u,
	}

	av, _ := dynamodbattribute.MarshalMap(feedback)

	av["createdTimestamp"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.Itoa(int(u))),
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamo.SonarFeedback),
	}

	suite.PutItem = input
}

func (suite *SubmitFeedbackTestSuite) TestSubmitFeedback_Success() {
	suite.DynamoMock.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	err := Handler(suite.Request)

	suite.Nil(err)
	suite.DynamoMock.AssertCalled(suite.T(), "PutItem", suite.PutItem)
}

func (suite *SubmitFeedbackTestSuite) TestSubmitFeedback_Fail() {
	suite.DynamoMock.On("PutItem", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(suite.Request)

	suite.NotNil(err)
	suite.DynamoMock.AssertCalled(suite.T(), "PutItem", suite.PutItem)
}

func (suite *SubmitFeedbackTestSuite) TestSubmitFeedback_FailEmail() {
	feedback := request.FeedbackRequest{
		Email:         nil,
		Activity:      "chatting",
		ActivityNotes: "wanted to chat but couldn't",
		Suggestion:    "make chat work",
	}

	suite.Request.Feedback = feedback

	err := Handler(suite.Request)

	suite.NotNil(err)
	suite.DynamoMock.AssertNotCalled(suite.T(), "PutItem", suite.PutItem)
}

type SendFeedbackTestSuite struct {
	suite.Suite
	SesMock                *mocks.SESAPI
	Input                  SendFeedbackInput
	SendTemplateEmailInput *ses.SendTemplatedEmailInput
}

func (suite *SendFeedbackTestSuite) SetupTest() {
	suite.SesMock = new(mocks.SESAPI)
	suite.Input = SendFeedbackInput{
		FeedbackData: []byte{},
		SesClient:    suite.SesMock,
		ConfigSet:    "test",
		Template:     "test",
		Domain:       "test",
	}
	sender := "support@" + suite.Input.Domain
	recipient := "madison@circulohealth.com"
	suite.SendTemplateEmailInput = &ses.SendTemplatedEmailInput{
		ConfigurationSetName: aws.String(suite.Input.ConfigSet),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Source:       aws.String(sender),
		Template:     aws.String(suite.Input.Template),
		TemplateData: aws.String(string(suite.Input.FeedbackData)),
	}
}

func (suite *SendFeedbackTestSuite) TestSendFeedback_Success() {
	suite.SesMock.On("SendTemplatedEmail", suite.SendTemplateEmailInput).Return(&ses.SendTemplatedEmailOutput{}, nil)

	err := SendFeedback(suite.Input)

	suite.Nil(err)
	suite.SesMock.AssertCalled(suite.T(), "SendTemplatedEmail", suite.SendTemplateEmailInput)
}

func (suite *SendFeedbackTestSuite) TestSendFeedback_Fail() {
	suite.SesMock.On("SendTemplatedEmail", suite.SendTemplateEmailInput).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := SendFeedback(suite.Input)

	suite.NotNil(err)
	suite.SesMock.AssertCalled(suite.T(), "SendTemplatedEmail", suite.SendTemplateEmailInput)
}

type MarshalFeedbackJSONTestSuite struct {
	suite.Suite
}

func (suite *MarshalFeedbackJSONTestSuite) TestSendFeedback_Success() {
	_, err := MarshalFeedbackJSON("Name", request.FeedbackRequest{
		Email:         new(string),
		UserName:      "test",
		Activity:      "test",
		ActivityNotes: "test",
		Suggestion:    "test",
	})

	suite.Nil(err)
}

/* Execute Suites */

func TestSubmitFeedbackTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitFeedbackTestSuite))
}

func TestSendFeedbackTestSuite(t *testing.T) {
	suite.Run(t, new(SendFeedbackTestSuite))
}

func TestMarshalFeedbackJSONTestSuite(t *testing.T) {
	suite.Run(t, new(MarshalFeedbackJSONTestSuite))
}
