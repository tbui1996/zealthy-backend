package iterator

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IteratorTestSuite struct {
	suite.Suite
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(message events.SQSMessage) error {
	args := m.Called(message)

	return args.Error(0)
}

func (s *IteratorTestSuite) TestHandle_CallsHandlerWithMessage() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	i := &Iterator{
		SQS: m,
	}

	handler := &MockHandler{}
	handler.On("Handle", mock.Anything).Return(nil)

	message := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	i.UseHandler(handler.Handle)

	// nolint errcheck
	i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message,
		},
	})

	handler.AssertCalled(s.T(), "Handle", message)
	handler.AssertNumberOfCalls(s.T(), "Handle", 1)
}

func (s *IteratorTestSuite) TestHandle_CallsHandlerWithMultipleMessages() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	i := &Iterator{
		SQS: m,
	}

	handler := &MockHandler{}
	handler.On("Handle", mock.Anything).Return(nil)

	message0 := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	message1 := events.SQSMessage{
		MessageId:     "2",
		ReceiptHandle: "2",
		Body:          "2",
	}

	i.UseHandler(handler.Handle)

	// nolint errcheck
	i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message0,
			message1,
		},
	})

	handler.AssertCalled(s.T(), "Handle", message0)
	handler.AssertCalled(s.T(), "Handle", message1)
	handler.AssertNumberOfCalls(s.T(), "Handle", 2)
}

func (s *IteratorTestSuite) TestHandle_CallsDeleteMessageWithSuccessfulMessages() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	handler := &MockHandler{}
	handler.On("Handle", mock.Anything).Return(nil)

	message := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	i.UseHandler(handler.Handle)

	// nolint errcheck
	i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message,
		},
	})

	expected := &sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: &message.ReceiptHandle,
	}

	m.AssertCalled(s.T(), "DeleteMessage", expected)
}

func (s *IteratorTestSuite) TestHandle_DoesNotCallsDeleteMessageWithFailedMessages() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	handler := &MockHandler{}
	handler.On("Handle", mock.Anything).Return(errors.New("test"))

	message := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	i.UseHandler(handler.Handle)

	// nolint errcheck
	i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message,
		},
	})

	m.AssertNumberOfCalls(s.T(), "DeleteMessage", 0)
}

func (s *IteratorTestSuite) TestHandle_CallsDeleteMessageWithSuccessfulMessagesAndDoesNotOnFailedMessages() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	message0 := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	message1 := events.SQSMessage{
		MessageId:     "2",
		ReceiptHandle: "2",
		Body:          "2",
	}

	handler := &MockHandler{}
	handler.On("Handle", message0).Return(nil)
	handler.On("Handle", message1).Return(errors.New("test"))

	i.UseHandler(handler.Handle)

	// nolint errcheck
	i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message0,
			message1,
		},
	})

	expected := &sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: &message0.ReceiptHandle,
	}

	handler.AssertCalled(s.T(), "Handle", message0)
	handler.AssertCalled(s.T(), "Handle", message1)
	m.AssertNumberOfCalls(s.T(), "DeleteMessage", 1)
	m.AssertCalled(s.T(), "DeleteMessage", expected)
}

func (s *IteratorTestSuite) TestHandle_ReturnsErrorIfMessageFails() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	message0 := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	message1 := events.SQSMessage{
		MessageId:     "2",
		ReceiptHandle: "2",
		Body:          "2",
	}

	expectedErrMessage := "test"
	handler := &MockHandler{}
	handler.On("Handle", message0).Return(nil)
	handler.On("Handle", message1).Return(errors.New(expectedErrMessage))

	i.UseHandler(handler.Handle)

	err := i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message0,
			message1,
		},
	})

	expected := &sqs.DeleteMessageInput{
		QueueUrl:      i.QueueUrl,
		ReceiptHandle: &message0.ReceiptHandle,
	}

	handler.AssertCalled(s.T(), "Handle", message0)
	handler.AssertCalled(s.T(), "Handle", message1)
	m.AssertNumberOfCalls(s.T(), "DeleteMessage", 1)
	m.AssertCalled(s.T(), "DeleteMessage", expected)
	s.NotNil(err)
	s.EqualError(err, expectedErrMessage)
}

func (s *IteratorTestSuite) TestHandle_ReturnsErrorMessagesIfMultipleMessagesFails() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	message0 := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	message1 := events.SQSMessage{
		MessageId:     "2",
		ReceiptHandle: "2",
		Body:          "2",
	}

	expectedErrMessage0 := "test0"
	expectedErrMessage1 := "test1"
	handler := &MockHandler{}
	handler.On("Handle", message0).Return(errors.New(expectedErrMessage0))
	handler.On("Handle", message1).Return(errors.New(expectedErrMessage1))

	i.UseHandler(handler.Handle)

	err := i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message0,
			message1,
		},
	})

	handler.AssertCalled(s.T(), "Handle", message0)
	handler.AssertCalled(s.T(), "Handle", message1)
	m.AssertNumberOfCalls(s.T(), "DeleteMessage", 0)
	s.NotNil(err)
	s.Contains(err.Error(), expectedErrMessage0)
	s.Contains(err.Error(), expectedErrMessage1)
}

func (s *IteratorTestSuite) TestHandle_ReturnsNilMessagesIfAllMessagesSucceed() {
	m := &mocks.SQSAPI{}
	m.On("DeleteMessage", mock.Anything).Return(nil, nil)

	// nolint goconst
	queueUrl := "test-queue-url"
	i := &Iterator{
		SQS:      m,
		QueueUrl: &queueUrl,
	}

	message0 := events.SQSMessage{
		MessageId:     "1",
		ReceiptHandle: "1",
		Body:          "1",
	}

	message1 := events.SQSMessage{
		MessageId:     "2",
		ReceiptHandle: "2",
		Body:          "2",
	}

	handler := &MockHandler{}
	handler.On("Handle", mock.Anything).Return(nil)

	i.UseHandler(handler.Handle)

	err := i.Handle(events.SQSEvent{
		Records: []events.SQSMessage{
			message0,
			message1,
		},
	})

	handler.AssertCalled(s.T(), "Handle", message0)
	handler.AssertCalled(s.T(), "Handle", message1)
	m.AssertNumberOfCalls(s.T(), "DeleteMessage", 2)
	s.Nil(err)
}

func TestIteratorTestSuite(t *testing.T) {
	suite.Run(t, new(IteratorTestSuite))
}
