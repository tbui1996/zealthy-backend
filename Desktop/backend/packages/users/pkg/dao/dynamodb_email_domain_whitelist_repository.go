package dao

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
)

// implements EmailDomainWhitelistRepository interface
type DynamoDBEmailDomainWhitelistRepository struct {
	DB      dynamodbiface.DynamoDBAPI
	Wrapper dynamo.Database
}

func NewDynamoDBEmailDomainWhitelistRepository(db dynamodbiface.DynamoDBAPI) *DynamoDBEmailDomainWhitelistRepository {
	dbWrapper := &dynamo.DynamoDatabase{
		TableName: dynamo.SonarEmailDomainWhitelist,
	}

	return &DynamoDBEmailDomainWhitelistRepository{
		DB:      db,
		Wrapper: dbWrapper,
	}
}

func (repo *DynamoDBEmailDomainWhitelistRepository) GetWhitelistDomain(domain string) (*model.EmailDomainWhitelist, *exception.SonarError) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(dynamo.SonarEmailDomainWhitelist),
		Key: map[string]*dynamodb.AttributeValue{
			"EmailDomain": {
				S: aws.String(domain),
			},
		},
	}

	result, err := repo.Wrapper.Get(input)

	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, err.Error())
	}

	if result.Item == nil || len(result.Item) == 0 {
		return nil, nil
	}

	var validDomain model.EmailDomainWhitelist
	err = dynamodbattribute.UnmarshalMap(result.Item, &validDomain)

	if err != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, err.Error())
	}

	if (model.EmailDomainWhitelist{}) == validDomain {
		return nil, nil
	}

	return &validDomain, nil
}
