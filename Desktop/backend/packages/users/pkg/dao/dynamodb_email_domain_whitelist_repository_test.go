package dao

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EmailDomainWhitelistRepositoryTestSuite struct {
	suite.Suite
}

var validEmailDomain = map[string]*dynamodb.AttributeValue{
	"EmailDomain": {
		S: aws.String("circulohealth.com"),
	},
}

var invalidDbResponse = map[string]*dynamodb.AttributeValue{
	"djnbvfjbsdfkbjn": {
		N: nil,
	},
}

func (suite *EmailDomainWhitelistRepositoryTestSuite) TestGetEntity_ShouldReturnWhiteListDomain() {
	mockEmailDomainDatabase := new(dynamo.MockDatabase)
	mockEmailDomainDatabase.On("Get", mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: validEmailDomain,
	}, nil)

	repo := DynamoDBEmailDomainWhitelistRepository{
		Wrapper: mockEmailDomainDatabase,
	}

	resp, err := repo.GetWhitelistDomain("circulohealth.com")

	theResp := *resp

	suite.Nil(err)
	suite.Equal("circulohealth.com", theResp.EmailDomain)
}

func (suite *EmailDomainWhitelistRepositoryTestSuite) TestGetEntity_ShouldNotReturnWhiteListDomain() {
	mockEmailDomainDatabase := new(dynamo.MockDatabase)
	mockEmailDomainDatabase.On("Get", mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{},
	}, nil)

	repo := DynamoDBEmailDomainWhitelistRepository{
		Wrapper: mockEmailDomainDatabase,
	}

	resp, err := repo.GetWhitelistDomain("nothing.com")

	suite.Nil(err)
	suite.Nil(resp)
}

func (suite *EmailDomainWhitelistRepositoryTestSuite) TestGetEntity_ShouldReturnDBError() {
	mockEmailDomainDatabase := new(dynamo.MockDatabase)
	mockEmailDomainDatabase.On("Get", mock.Anything).Return(nil, errors.New("could not scan DB"))

	repo := DynamoDBEmailDomainWhitelistRepository{
		Wrapper: mockEmailDomainDatabase,
	}

	resp, err := repo.GetWhitelistDomain("noworries.com")

	suite.NotNil(err)
	suite.Equal(400, err.StatusCode)
	suite.Nil(resp)
}

func (suite *EmailDomainWhitelistRepositoryTestSuite) TestGetEntity_ShouldReturnNoItemsInvalidStruct() {
	mockEmailDomainDatabase := new(dynamo.MockDatabase)
	mockEmailDomainDatabase.On("Get", mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: invalidDbResponse,
	}, nil)

	repo := DynamoDBEmailDomainWhitelistRepository{
		Wrapper: mockEmailDomainDatabase,
	}

	resp, err := repo.GetWhitelistDomain("what.com")

	suite.Nil(err)
	suite.Nil(resp)
}

/* Execute Suites */
func TestDynamoDBChatSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EmailDomainWhitelistRepositoryTestSuite))
}

func TestNewDynamoDBEmailDomainWhitelistRepository(t *testing.T) {
	mockDB := new(mocks.DynamoDBAPI)
	wrapper := NewDynamoDBEmailDomainWhitelistRepository(mockDB)
	assert.Equal(t, wrapper.DB, mockDB)
}
