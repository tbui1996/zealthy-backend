package connection

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type OnlineCheckerSuite struct {
	suite.Suite
}

func (s *OnlineCheckerSuite) Test_Success() {
	mockDb := new(mocks.DynamoDBAPI)
	logger := zaptest.NewLogger(s.T())
	checkerOne := OnlineChecker{
		TableName: "1",
		DB:        mockDb,
		Logger:    logger,
	}

	checkerTwo := OnlineChecker{
		TableName: "2",
		DB:        mockDb,
		Logger:    logger,
	}

	mockDb.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{Count: aws.Int64(1)}, nil)

	_, _ = checkerOne.IsUserOnline("123")
	_, _ = checkerTwo.IsUserOnline("321")

	mockDb.AssertCalled(s.T(), "Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
		return *input.TableName == "1" && *input.ExpressionAttributeValues[":userId"].S == "123"
	}))

	mockDb.AssertCalled(s.T(), "Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
		return *input.TableName == "2" && *input.ExpressionAttributeValues[":userId"].S == "321"
	}))
}

func (s *OnlineCheckerSuite) Test_Fail() {
	mockDb := new(mocks.DynamoDBAPI)
	logger := zaptest.NewLogger(s.T())
	checkerOne := OnlineChecker{
		TableName: "1",
		DB:        mockDb,
		Logger:    logger,
	}

	mockDb.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, fmt.Errorf("uh-oh"))

	dto, err := checkerOne.IsUserOnline("123")

	s.False(dto.IsOnline)
	s.NotNil(err)
}

func (s *OnlineCheckerSuite) Test_IsOnline() {
	mockDb := new(mocks.DynamoDBAPI)
	logger := zaptest.NewLogger(s.T())
	checkerOne := OnlineChecker{
		TableName: "1",
		DB:        mockDb,
		Logger:    logger,
	}

	mockDb.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{Count: aws.Int64(1)}, nil)

	dto, err := checkerOne.IsUserOnline("123")

	s.True(dto.IsOnline)
	s.Equal(dto.UserId, "123")
	s.Nil(err)
}

func (s *OnlineCheckerSuite) Test_IsNotOnline() {
	mockDb := new(mocks.DynamoDBAPI)
	logger := zaptest.NewLogger(s.T())
	checkerOne := OnlineChecker{
		TableName: "1",
		DB:        mockDb,
		Logger:    logger,
	}

	mockDb.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{Count: aws.Int64(0)}, nil)

	dto, err := checkerOne.IsUserOnline("123")

	s.False(dto.IsOnline)
	s.Equal(dto.UserId, "123")
	s.Nil(err)
}

func TestOnlineChecker(t *testing.T) {
	suite.Run(t, new(OnlineCheckerSuite))
}
