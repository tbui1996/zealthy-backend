package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/chatHelper"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/requestConfig"
	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/support/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

func CreateDefaultEvent() events.APIGatewayWebsocketProxyRequest {
	return CreateEvent("chat", "hi")
}

func CreateEvent(payloadType string, message string) events.APIGatewayWebsocketProxyRequest {
	request := request.SupportRequestSend{
		Action: "whatever",
		Payload: request.Payload{
			Type:    payloadType,
			Message: message,
		},
	}
	deserializedRequest, _ := json.Marshal(request)
	event := events.APIGatewayWebsocketProxyRequest{
		Body: string(deserializedRequest),
	}

	return event
}

type HandlerSuite struct {
	suite.Suite
	OgHandleChat func(
		config *requestConfig.APIGatewayWebsocketProxyRequest,
		message string,
		repo iface.ChatSessionRepository,
		client *router.Session,
	) error
	OGHandleReceipt func(payload string, time int64, repo iface.ChatSessionRepository) error
	OGHandleTyping  func(
		config *requestConfig.APIGatewayWebsocketProxyRequest,
		message string,
		repo iface.ChatSessionRepository,
		client *router.Session,
	) error
}

type ChatEventMock struct {
	mock.Mock
}

func (m *ChatEventMock) MockHandleChat(
	config *requestConfig.APIGatewayWebsocketProxyRequest,
	message string,
	repo iface.ChatSessionRepository,
	client *router.Session,
) error {
	args := m.Called(config, message, repo, client)

	return args.Error(0)
}

func (m *ChatEventMock) MockHandleReceipt(payload string, time int64, repo iface.ChatSessionRepository) error {
	args := m.Called(payload, time, repo)

	return args.Error(0)
}

func (m *ChatEventMock) MockHandleTyping(
	config *requestConfig.APIGatewayWebsocketProxyRequest,
	message string,
	repo iface.ChatSessionRepository,
	client *router.Session,
) error {
	args := m.Called(config, message, repo, client)

	return args.Error(0)
}

var chatEventMock *ChatEventMock

func (s *HandlerSuite) SetupSuite() {
	s.OgHandleChat = HandleChat
	s.OGHandleReceipt = chatHelper.HandleReadReceipt
	s.OGHandleTyping = HandleTyping
}

func (s *HandlerSuite) BeforeTest(suiteName, testName string) {
	chatEventMock = &ChatEventMock{}
	HandleChat = chatEventMock.MockHandleChat
	chatHelper.HandleReadReceipt = chatEventMock.MockHandleReceipt
	HandleTyping = chatEventMock.MockHandleTyping
}

func (s *HandlerSuite) Teardown(t *testing.T) {
	HandleChat = s.OgHandleChat
	chatHelper.HandleReadReceipt = s.OGHandleReceipt
	HandleTyping = s.OGHandleTyping
}

func (s *HandlerSuite) Test_HandlesBadBody() {
	event := CreateDefaultEvent()
	event.Body = "bad json"
	config := requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  event,
	}
	repo := new(mocks.ChatSessionRepository)

	_, err := Handler(&config, repo, &router.Session{})

	s.Equal(http.StatusInternalServerError, err.StatusCode)
	s.Equal("unable to unmarshal body: bad json", err.Error())
}

func (s *HandlerSuite) Test_HandlesSuccessfulChat() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("chat", "hello!"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}

	chatEventMock.On("MockHandleChat", mock.Anything, "hello!", repo, client).Return(nil)

	result, _ := Handler(config, repo, client)

	s.Equal("OK", result.Body)
	s.Equal(http.StatusOK, result.StatusCode)
}

func (s *HandlerSuite) Test_HandlesFailedChat() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("chat", "hello!"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}

	chatEventMock.On("MockHandleChat", mock.Anything, "hello!", repo, client).Return(fmt.Errorf("uh ohhhhh"))

	_, err := Handler(config, repo, client)

	s.Equal("uh ohhhhh", err.Error())
	s.Equal(http.StatusInternalServerError, err.StatusCode)
}

func (s *HandlerSuite) Test_HandlesReadReceiptSuccess() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("read_receipt", "hello!"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}
	chatEventMock.On("MockHandleReceipt", "hello!", mock.Anything, repo).Return(nil)

	result, _ := Handler(config, repo, client)

	s.Equal("OK", result.Body)
	s.Equal(http.StatusOK, result.StatusCode)
}

func (s *HandlerSuite) Test_HandlesReadReceiptFailure() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("read_receipt", "hello!"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}
	chatEventMock.On("MockHandleReceipt", "hello!", mock.Anything, repo).Return(fmt.Errorf("whyyyyy"))

	_, err := Handler(config, repo, client)

	s.Equal("whyyyyy", err.Error())
	s.Equal(http.StatusInternalServerError, err.StatusCode)
}

func (s *HandlerSuite) Test_HandlesTypingSuccess() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("typing", "Fake Payload"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}
	chatEventMock.On("MockHandleTyping", mock.Anything, "Fake Payload", repo, client).Return(nil)

	result, _ := Handler(config, repo, client)

	s.Equal("OK", result.Body)
	s.Equal(http.StatusOK, result.StatusCode)
}

func (s *HandlerSuite) Test_HandlesTypingFailure() {
	config := &requestConfig.APIGatewayWebsocketProxyRequest{
		Logger: zaptest.NewLogger(s.T()),
		Event:  CreateEvent("typing", "Fake Payload"),
	}
	repo := new(mocks.ChatSessionRepository)
	client := &router.Session{}
	chatEventMock.On("MockHandleTyping", mock.Anything, "Fake Payload", repo, client).Return(errors.New("FAKE ERROR"))

	_, err := Handler(config, repo, client)

	s.Equal("FAKE ERROR", err.Error())
	s.Equal(http.StatusBadRequest, err.StatusCode)
}

func TestSendHandler(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
