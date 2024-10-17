package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FileDownloadTestSuite struct {
	suite.Suite
	Request FileDownloadRequest
	S3Mock  *mocks.S3API
	Result  *request.Request
}

func (suite *FileDownloadTestSuite) SetupTest() {
	svc := new(mocks.S3API)
	suite.S3Mock = svc

	suite.Result = &request.Request{
		Operation: &request.Operation{
			BeforePresignFn: nil,
		},
		Handlers: request.Handlers{
			Validate:         request.HandlerList{AfterEachFn: nil},
			Build:            request.HandlerList{AfterEachFn: nil},
			BuildStream:      request.HandlerList{AfterEachFn: nil},
			Sign:             request.HandlerList{AfterEachFn: nil},
			Send:             request.HandlerList{AfterEachFn: nil},
			ValidateResponse: request.HandlerList{AfterEachFn: nil},
			Unmarshal:        request.HandlerList{AfterEachFn: nil},
			UnmarshalStream:  request.HandlerList{AfterEachFn: nil},
			UnmarshalMeta:    request.HandlerList{AfterEachFn: nil},
			UnmarshalError:   request.HandlerList{AfterEachFn: nil},
			Retry:            request.HandlerList{AfterEachFn: nil},
			AfterRetry:       request.HandlerList{AfterEachFn: nil},
			CompleteAttempt:  request.HandlerList{AfterEachFn: nil},
			Complete:         request.HandlerList{AfterEachFn: nil},
		},
		NotHoist:   false,
		ExpireTime: 5 * time.Minute,
		Error:      nil,
		HTTPRequest: &http.Request{
			Host: "",
			URL: &url.URL{
				Host: "google.com",
			},
		},
		SignedHeaderVals: http.Header{},
	}

	suite.Request = FileDownloadRequest{
		S3:         svc,
		Logger:     zaptest.NewLogger(suite.T()),
		BucketName: "test",
		FileId:     "1",
	}
}

func (suite *FileDownloadTestSuite) TestFileDownload_Success() {
	suite.S3Mock.On("GetObjectRequest", mock.Anything).Return(suite.Result, nil)

	res, err := handler(suite.Request)

	suite.Nil(err)
	suite.NotNil(res)
	suite.S3Mock.AssertCalled(suite.T(), "GetObjectRequest", mock.Anything)
}

func (suite *FileDownloadTestSuite) TestFileDownload_FailPresign() {
	res := suite.Result
	res.Error = errors.New("FAKE TEST ERROR, IGNORE")

	suite.S3Mock.On("GetObjectRequest", mock.Anything).Return(res, nil)

	r, err := handler(suite.Request)

	suite.NotNil(err)
	suite.Nil(r)
	suite.S3Mock.AssertCalled(suite.T(), "GetObjectRequest", mock.Anything)
}

func TestFileDownloadTestSuite(t *testing.T) {
	suite.Run(t, new(FileDownloadTestSuite))
}
