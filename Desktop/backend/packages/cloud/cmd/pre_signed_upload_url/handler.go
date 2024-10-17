package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	data "github.com/circulohealth/sonar-backend/packages/cloud/mime"
	"go.uber.org/zap"
	"strings"
	"time"
)

type PreSignedUploadUrlRequest struct {
	S3API      s3iface.S3API
	Logger     *zap.Logger
	BucketName string
	UniqueKey  string
	Filename   string
}

func Handler(req PreSignedUploadUrlRequest) (*string, error) {
	fileParts := strings.Split(req.Filename, ".")
	fileExt := fmt.Sprintf(".%s", fileParts[len(fileParts)-1])
	contentType, ok := data.MimeTypes[fileExt]

	if !ok {
		errMsg := fmt.Errorf("expected a valid content type for (%s)", fileExt)
		req.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	r, _ := req.S3API.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(req.BucketName),
		Key:         aws.String(req.UniqueKey),
		ContentType: aws.String(contentType),
		Metadata:    map[string]*string{"FileName": &req.Filename},
	})

	req.Logger.Debug(fmt.Sprintf("generating pre signed url for object key %s, bucket %s", req.UniqueKey, req.BucketName))

	// nolint gomnd
	str, err := r.Presign(5 * time.Minute)

	if err != nil {
		errMsg := fmt.Errorf("unable to generate presigned url for key %s, bucket %s. err (%s)", req.UniqueKey, req.BucketName, err)
		req.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	req.Logger.Debug("generated pre signed url")
	return &str, nil
}
