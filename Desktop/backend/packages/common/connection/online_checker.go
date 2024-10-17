package connection

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/connection/dto"
	"go.uber.org/zap"
)

type OnlineChecker struct {
	TableName string
	DB        dynamodbiface.DynamoDBAPI
	Logger    *zap.Logger
}

func (checker *OnlineChecker) IsUserOnline(userID string) (dto.UserOnlineStatus, error) {
	query := &dynamodb.QueryInput{
		TableName:              aws.String(checker.TableName),
		KeyConditionExpression: aws.String("UserID = :userId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {
				S: aws.String(userID),
			},
		},
	}

	res, err := checker.DB.Query(query)
	if err != nil {
		checker.Logger.Error(err.Error())
		return dto.UserOnlineStatus{}, err
	}
	count := *res.Count

	return dto.UserOnlineStatus{IsOnline: count > 0, UserId: userID}, nil
}
