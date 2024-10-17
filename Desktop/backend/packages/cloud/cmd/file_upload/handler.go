package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	data "github.com/circulohealth/sonar-backend/packages/cloud/mime"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/request"
	"strconv"
	"strings"
	"time"

	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FileUploadHandler struct {
	S3            s3iface.S3API
	DB            *gorm.DB
	BucketName    string
	Logger        *zap.Logger
	Username      string
	UploadRequest request.FileUploadRequest
}

func handler(req *FileUploadHandler) ([]byte, error) {
	filename := req.UploadRequest.Filename
	fileId := req.UploadRequest.FileId
	fileParts := strings.Split(filename, ".")
	fileExt := fmt.Sprintf(".%s", fileParts[len(fileParts)-1])
	fileMimeType, ok := data.MimeTypes[fileExt]
	if !ok {
		errMsg := fmt.Errorf("expected a valid content type for (%s)", fileExt)
		req.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	currTime := time.Now()
	file := model.File{
		FileID:           fileId,
		FileName:         filename,
		FileMimetype:     fileMimeType,
		SendUserID:       req.Username,
		ChatID:           req.UploadRequest.ChatId,
		DateUploaded:     currTime,
		DateLastAccessed: currTime,
		FilePath:         fmt.Sprintf("https://s3.us-east-2.amazonaws.com/%s/%s", req.BucketName, fileId),
		MemberID:         "",
		DeletedAt:        nil,
	}
	result := req.DB.Create(&file)

	if result.Error != nil {
		_, err := req.S3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(req.BucketName),
			Key:    aws.String(fileId),
		})
		errMsg := fmt.Errorf("unable to create file in database: (%s)", result.Error.Error())
		if err != nil {
			errMsg = fmt.Errorf("unable to delete object with key (%s): (%s)", fileId, err.Error())
			req.Logger.Error(errMsg.Error())
			return nil, errMsg
		}

		req.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	resId := strconv.Itoa(file.ID)
	req.Logger = req.Logger.With(zap.String("fileID", resId))

	fileResponse, err := json.Marshal(model.FileUploadResponse{
		FileID: resId,
	})

	if err != nil || resId == "0" {
		errMsg := fmt.Errorf("unable to marshal file response (%s)", err)
		req.Logger.Error(errMsg.Error())
		return nil, errMsg
	}

	return fileResponse, nil
}
