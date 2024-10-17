package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/logging"
	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func CreateDefaultSqsMessage() events.SQSMessage {
	return CreateSqsMessage("read_receipt", "hi")
}

func CreateSqsMessage(messageType string, message string) events.SQSMessage {
	loggerFields := logging.LoggerFields{}
	loggerFieldsString, _ := loggerFields.ToString()
	connectionId := "ConnectionId"
	serializedBody := &request.SupportRequestReceive{
		Type:    messageType,
		Message: message,
	}
	deserializedBody, _ := json.Marshal(serializedBody)
	return events.SQSMessage{
		MessageId:              "",
		ReceiptHandle:          "",
		Body:                   string(deserializedBody),
		Md5OfBody:              "",
		Md5OfMessageAttributes: "",
		Attributes:             map[string]string{},
		MessageAttributes: map[string]events.SQSMessageAttribute{
			"LoggerFields": {StringValue: &loggerFieldsString},
			"ConnectionId": {StringValue: &connectionId},
		},
		EventSourceARN: "",
		EventSource:    "",
		AWSRegion:      "",
	}
}

type HandleReadReceiptMock struct {
	mock.Mock
}

func (m *HandleReadReceiptMock) MockImpl(payload string, time int64, repo iface.ChatSessionRepository) error {
	args := m.Called(payload, time, repo)

	return args.Error(0)
}

var mockHandleReceipt *HandleReadReceiptMock

func DefaultMockHandleReceipt() {
	mockHandleReceipt.On("MockImpl", mock.Anything, mock.Anything, mock.Anything).Return(nil)
}

type HandlerSuite struct {
	suite.Suite
	OGHandleReceipt func(payload string, time int64, repo iface.ChatSessionRepository) error
}

func (s *HandlerSuite) SetupSuite() {
	s.OGHandleReceipt = chatHelper.HandleReadReceipt
}

func (s *HandlerSuite) BeforeTest(suiteName, testName string) {
	mockHandleReceipt = &HandleReadReceiptMock{}
	chatHelper.HandleReadReceipt = mockHandleReceipt.MockImpl
}

func (s *HandlerSuite) Teardown(t *testing.T) {
	chatHelper.HandleReadReceipt = s.OGHandleReceipt
}

func (s *HandlerSuite) Test_DoesNotPanicMissingConnectionId() {
	DefaultMockHandleReceipt()
	repo := new(mocks.ChatSessionRepository)
	message := CreateDefaultSqsMessage()
	delete(message.MessageAttributes, "ConnectionId")

	err := Handler(message, repo)

	s.Nil(err)
}

func (s *HandlerSuite) Test_HandlesChatHelperError() {
	repo := new(mocks.ChatSessionRepository)
	message := CreateDefaultSqsMessage()

	mockHandleReceipt.On("MockImpl", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("lol wat"))

	err := Handler(message, repo)

	s.Equal("lol wat", err.Error())
}

func (s *HandlerSuite) Test_SendsCorrectParametersToChatHelper() {
	DefaultMockHandleReceipt()
	repo := new(mocks.ChatSessionRepository)
	message := CreateSqsMessage("read_receipt", "heyyyy!")

	err := Handler(message, repo)

	s.Nil(err)
	mockHandleReceipt.AssertCalled(s.T(), "MockImpl", "heyyyy!", mock.AnythingOfType("int64"), repo)
}

func TestRecieveHandler(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
