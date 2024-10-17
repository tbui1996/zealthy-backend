package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"go.uber.org/zap"
	"time"
)

type FileDownloadRequest struct {
	S3         s3iface.S3API
	BucketName string
	FileId     string
	Logger     *zap.Logger
}

func handler(req FileDownloadRequest) (*string, error) {
	r, _ := req.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(req.BucketName),
		Key:    aws.String(req.FileId),
	})

	// nolint gomnd
	str, err := r.Presign(168 * time.Hour)
	if err != nil {
		return nil, err
	}

	return &str, nil
}
