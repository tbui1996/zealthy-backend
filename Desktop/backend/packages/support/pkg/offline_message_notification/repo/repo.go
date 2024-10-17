package repo

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/constants"
	"go.uber.org/zap"
)

type OfflineMessageNotificationRepo struct {
	DynamoDB dynamodbiface.DynamoDBAPI
	Logger   *zap.Logger
}

const DECIMAL_BASE int = 10

func epoch() string {
	return strconv.FormatInt(time.Now().Unix(), DECIMAL_BASE)
}

func NewOfflineMessageNotificationRepoWithLogger(logger *zap.Logger) *OfflineMessageNotificationRepo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	return &OfflineMessageNotificationRepo{
		DynamoDB: db,
		Logger:   logger,
	}
}

func (repo *OfflineMessageNotificationRepo) Create(userID string) (bool, error) {
	putItemRequest := dynamodb.PutItemInput{
		TableName: aws.String(dynamo.OfflineMessageNotifications),
		Item: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
			"Status": {
				S: aws.String(constants.PENDING_NOTIFICATION),
			},
			"CreatedAt": {
				N: aws.String(epoch()),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(UserID)"), // this makes it so only one notificatio can be pending at a time and the timer doesn't restart
	}

	_, err := repo.DynamoDB.PutItem(&putItemRequest)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				repo.Logger.Info("notification was not created because it already existed, skipping..")
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (repo *OfflineMessageNotificationRepo) Remove(userID string) error {
	request := &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamo.OfflineMessageNotifications),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	}

	_, err := repo.DynamoDB.DeleteItem(request)

	return err
}

func (repo *OfflineMessageNotificationRepo) UpdateStatus(userId string, status string) (bool, error) {
	update := &dynamodb.UpdateItemInput{
		TableName:        aws.String(dynamo.OfflineMessageNotifications),
		UpdateExpression: aws.String("set #status = :s, SentAt = :now"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userId),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(status),
			},
			":now": {
				N: aws.String(epoch()),
			},
		},
		// this makes sure the user didn't come online, notification records are deleted when the user comes online
		ConditionExpression: aws.String("attribute_exists(UserID)"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
	}

	_, err := repo.DynamoDB.UpdateItem(update)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				repo.Logger.Info("notification already sent or the user came online")
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
