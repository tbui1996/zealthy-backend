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

type PreSignedUploadUrlTestSuite struct {
	suite.Suite
	Request PreSignedUploadUrlRequest
	S3Mock  *mocks.S3API
	Result  *request.Request
}

func (suite *PreSignedUploadUrlTestSuite) SetupTest() {
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

	suite.Request = PreSignedUploadUrlRequest{
		S3API:      svc,
		Logger:     zaptest.NewLogger(suite.T()),
		BucketName: "test",
		UniqueKey:  "test",
		Filename:   "test.png",
	}
}

func (suite *PreSignedUploadUrlTestSuite) TestPreSignedUploadUrl_Success() {
	suite.S3Mock.On("PutObjectRequest", mock.Anything).Return(suite.Result, nil)

	res, err := Handler(suite.Request)

	suite.Nil(err)
	suite.NotNil(res)
	suite.S3Mock.AssertCalled(suite.T(), "PutObjectRequest", mock.Anything)
}

func (suite *PreSignedUploadUrlTestSuite) TestPreSignedUploadUrl_FailFileType() {
	req := suite.Request
	req.Filename = "bunk.sdfgdsfgsdfg"

	res, err := Handler(req)

	suite.NotNil(err)
	suite.Nil(res)
	suite.S3Mock.AssertNotCalled(suite.T(), "PutObjectRequest", mock.Anything)
}

func (suite *PreSignedUploadUrlTestSuite) TestPreSignedUploadUrl_FailPresign() {
	res := suite.Result
	res.Error = errors.New("FAKE TEST ERROR, IGNORE")

	suite.S3Mock.On("PutObjectRequest", mock.Anything).Return(res, nil)

	r, err := Handler(suite.Request)

	suite.NotNil(err)
	suite.Nil(r)
	suite.S3Mock.AssertCalled(suite.T(), "PutObjectRequest", mock.Anything)
}

func TestPreSignedUploadUrlTestSuite(t *testing.T) {
	suite.Run(t, new(PreSignedUploadUrlTestSuite))
}
