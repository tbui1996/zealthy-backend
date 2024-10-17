package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/router/pkg/request"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceReceiveTestSuite struct {
	suite.Suite
}

func getEvent() events.APIGatewayWebsocketProxyRequest {
	b, _ := json.Marshal(request.Payload{
		Action:  "action",
		Payload: "payload",
	})

	return events.APIGatewayWebsocketProxyRequest{
		Body: string(b),
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			Authorizer: map[string]interface{}{},
			RouteKey:   "route",
			RequestID:  "request",
			Identity: events.APIGatewayRequestIdentity{
				SourceIP:  "0.0.0.0",
				UserAgent: "Mozilla",
			},
			ConnectionID: "1",
		},
	}
}

func (suite *ServiceReceiveTestSuite) TestServiceReceive_Success() {
	mockSQS := new(mocks.SQSAPI)

	url := "url"
	mockSQS.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: &url,
	}, nil)

	mockSQS.On("SendMessage", mock.Anything).Return(nil, nil)

	err := Handler(ServiceReceiveRequest{
		Name:             "test",
		ReceiveQueueName: "receive",
		SendQueueName:    "send",
		SQS:              mockSQS,
		Event:            getEvent(),
		Logger:           zaptest.NewLogger(suite.T()),
	})

	suite.Nil(err)
	mockSQS.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSQS.AssertCalled(suite.T(), "SendMessage", mock.Anything)
}

func (suite *ServiceReceiveTestSuite) TestServiceReceive_FailOnGetUrl() {
	mockSQS := new(mocks.SQSAPI)

	mockSQS.On("GetQueueUrl", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(ServiceReceiveRequest{
		Name:             "test",
		ReceiveQueueName: "receive",
		SendQueueName:    "send",
		SQS:              mockSQS,
		Event:            getEvent(),
		Logger:           zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockSQS.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSQS.AssertNotCalled(suite.T(), "SendMessage", mock.Anything)
}

func (suite *ServiceReceiveTestSuite) TestServiceReceive_FailOnSendMessage() {
	mockSQS := new(mocks.SQSAPI)

	url := "url"
	mockSQS.On("GetQueueUrl", mock.Anything).Return(&sqs.GetQueueUrlOutput{
		QueueUrl: &url,
	}, nil)

	mockSQS.On("SendMessage", mock.Anything).Return(nil, errors.New("FAKE TEST ERROR, IGNORE"))

	err := Handler(ServiceReceiveRequest{
		Name:             "test",
		ReceiveQueueName: "receive",
		SendQueueName:    "send",
		SQS:              mockSQS,
		Event:            getEvent(),
		Logger:           zaptest.NewLogger(suite.T()),
	})

	suite.NotNil(err)
	mockSQS.AssertCalled(suite.T(), "GetQueueUrl", mock.Anything)
	mockSQS.AssertCalled(suite.T(), "SendMessage", mock.Anything)
}

func TestServiceReceiveTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceReceiveTestSuite))
}
