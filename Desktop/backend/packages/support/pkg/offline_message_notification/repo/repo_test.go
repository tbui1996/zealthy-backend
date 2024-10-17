package repo

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/offline_message_notification/constants"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

/* Create Tests Start */
type CreateSuite struct {
	suite.Suite
	logger *zap.Logger
	db     *mocks.DynamoDBAPI
}

const TEST_USER_ID string = "test-user-id"

func (s *CreateSuite) SetupTest() {

	s.logger = zaptest.NewLogger(s.T())
	s.db = new(mocks.DynamoDBAPI)
}

func (s *CreateSuite) TestCallsDynamoCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("PutItem", mock.Anything).Return(nil, nil)

	created, err := repo.Create(TEST_USER_ID)

	s.db.AssertCalled(s.T(), "PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
		return *input.Item["UserID"].S == TEST_USER_ID &&
			*input.Item["Status"].S == constants.PENDING_NOTIFICATION &&
			*input.ConditionExpression == "attribute_not_exists(UserID)" &&
			*input.TableName == dynamo.OfflineMessageNotifications
	}))

	s.Nil(err)
	s.True(created)
}

func (s *CreateSuite) TestHandlesAwsErrorCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("PutItem", mock.Anything).Return(nil, awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "oops", nil))

	created, err := repo.Create(TEST_USER_ID)

	s.Nil(err)
	s.False(created)
}

func (s *CreateSuite) TestHandlesOtherErrorCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("PutItem", mock.Anything).Return(nil, fmt.Errorf("oh noes"))

	created, err := repo.Create(TEST_USER_ID)

	s.NotNil(err)
	s.False(created)
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

/* Create Tests End */

/* Remove Tests Start */
type RemoveSuite struct {
	suite.Suite
	logger *zap.Logger
	db     *mocks.DynamoDBAPI
}

func (s *RemoveSuite) SetupTest() {
	s.logger = zaptest.NewLogger(s.T())
	s.db = new(mocks.DynamoDBAPI)
}

/* Create Tests End */

/* Remove Tests Start */
func (s *RemoveSuite) TestSuccess() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("DeleteItem", mock.Anything).Return(nil, nil)

	err := repo.Remove(TEST_USER_ID)

	s.Nil(err)
}

func (s *RemoveSuite) TestError() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("DeleteItem", mock.Anything).Return(nil, fmt.Errorf("Oh Noes!!"))

	err := repo.Remove(TEST_USER_ID)

	s.NotNil(err)
}

func TestRemove(t *testing.T) {
	suite.Run(t, new(RemoveSuite))
}

/* Remove Tests End */

/* Create Tests Start */
type UpdateStatusSuite struct {
	suite.Suite
	item   map[string]*dynamodb.AttributeValue
	logger *zap.Logger
	db     *mocks.DynamoDBAPI
}

func (s *UpdateStatusSuite) SetupTest() {
	s.item = map[string]*dynamodb.AttributeValue{
		"UserID": {
			S: aws.String(TEST_USER_ID),
		},
		"Status": {
			S: aws.String(constants.PENDING_NOTIFICATION),
		},
		"CreatedAt": {
			N: aws.String(epoch()),
		},
	}

	s.logger = zaptest.NewLogger(s.T())
	s.db = new(mocks.DynamoDBAPI)
}

func (s *UpdateStatusSuite) TestCallsDynamoCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("UpdateItem", mock.Anything).Return(nil, nil)

	updated, err := repo.UpdateStatus(TEST_USER_ID, "LOLOL")

	// update is correct
	s.db.AssertCalled(s.T(), "UpdateItem", mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
		return *input.Key["UserID"].S == TEST_USER_ID &&
			*input.UpdateExpression == "set #status = :s, SentAt = :now" &&
			*input.ExpressionAttributeValues[":s"].S == "LOLOL" &&
			*input.ExpressionAttributeNames["#status"] == "Status" &&
			*input.ConditionExpression == "attribute_exists(UserID)"
	}))

	s.Nil(err)
	s.True(updated)
}

func (s *UpdateStatusSuite) TestHandlesAwsErrorCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("UpdateItem", mock.Anything).Return(nil, awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "oops", nil))

	updated, err := repo.UpdateStatus(TEST_USER_ID, "some-status")

	s.Nil(err)
	s.False(updated)
}

func (s *UpdateStatusSuite) TestHandlesOtherErrorCorrectly() {
	repo := OfflineMessageNotificationRepo{
		s.db,
		s.logger,
	}

	s.db.On("UpdateItem", mock.Anything).Return(nil, fmt.Errorf("oh noes"))

	updated, err := repo.UpdateStatus(TEST_USER_ID, "wat")

	s.NotNil(err)
	s.False(updated)
}

func TestUpdateStatus(t *testing.T) {
	suite.Run(t, new(UpdateStatusSuite))
}

/* Create Tests End */
