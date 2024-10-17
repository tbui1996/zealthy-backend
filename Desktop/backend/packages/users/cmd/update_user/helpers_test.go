package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type HandleOrganizationUpdateSuite struct {
	suite.Suite
	externalUserOrganizationMapper *iface.MockExternalUserOrganization
	registry                       *mapper.MockRegistryAPI
	userToUpdate                   *model.ExternalUser
	organization                   *model.ExternalUserOrganization
}

func (suite *HandleOrganizationUpdateSuite) SetupTest() {
	suite.registry = new(mapper.MockRegistryAPI)
	suite.externalUserOrganizationMapper = new(iface.MockExternalUserOrganization)
	suite.organization = &model.ExternalUserOrganization{
		ID: 1,
	}
	suite.userToUpdate = &model.ExternalUser{
		ID:       "1",
		Username: "Username",
		Email:    "name@circulohealth.com",
		Status:   "",
	}
}

func (suite *HandleOrganizationUpdateSuite) TestUpdateUserOrganization_Success() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Find", 1).Return(suite.organization, nil)

	actualErr := handleOrganizationUpdate(1, UpdateOrganizationDependencies{
		User:     suite.userToUpdate,
		Registry: suite.registry,
	})

	suite.NoError(actualErr)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Find", 1)
}

func (suite *HandleOrganizationUpdateSuite) TestUpdateUserOrganization_Fail() {
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Find", 1).Return(nil, errors.New("FAKE ERROR"))

	actualErr := handleOrganizationUpdate(1, UpdateOrganizationDependencies{
		User:     suite.userToUpdate,
		Registry: suite.registry,
	})

	suite.Error(actualErr)
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Find", 1)
}

func (suite *HandleOrganizationUpdateSuite) TestUpdateUserOrganizationToNil_Success() {
	actualErr := handleOrganizationUpdate(0, UpdateOrganizationDependencies{
		User:     suite.userToUpdate,
		Registry: suite.registry,
	})

	suite.Nil(actualErr)
	suite.externalUserOrganizationMapper.AssertNotCalled(suite.T(), "Find", 1)
}

type SendConfirmationNotificationSuite struct {
	suite.Suite
	userToUpdate     *model.ExternalUser
	dynamoQueryInput *dynamodb.QueryInput
	dependencies     GetSonarWebsocketConnectionDependencies
}

func (suite *SendConfirmationNotificationSuite) SetupTest() {
	suite.dependencies = GetSonarWebsocketConnectionDependencies{
		Logger: zaptest.NewLogger(suite.T()),
		Db:     new(mocks.DynamoDBAPI),
	}

	suite.userToUpdate = &model.ExternalUser{
		ID:       "1",
		Username: "Username",
		Email:    "name@circulohealth.com",
		Status:   "",
	}
	suite.dynamoQueryInput = &dynamodb.QueryInput{
		TableName:              aws.String(dynamo.SonarUnconfirmedWebsocketConnections),
		KeyConditionExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(suite.userToUpdate.Email),
			},
		},
		ProjectionExpression: aws.String("ConnectionId"),
	}
}

func (suite *SendConfirmationNotificationSuite) TestSendConfirmation_Success() {
	suite.dependencies.Db.(*mocks.DynamoDBAPI).On("Query", suite.dynamoQueryInput).Return(&dynamodb.QueryOutput{}, nil)

	actual, actualErr := getSonarWebsocketConnection("name@circulohealth.com", suite.dependencies)

	suite.NotNil(actual)
	suite.NoError(actualErr)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func (suite *SendConfirmationNotificationSuite) TestSendConfirmation_Fail() {
	suite.dependencies.Db.(*mocks.DynamoDBAPI).On("Query", suite.dynamoQueryInput).Return(&dynamodb.QueryOutput{}, errors.New("FAKE ERROR"))

	actual, actualErr := getSonarWebsocketConnection("name@circulohealth.com", suite.dependencies)

	suite.Nil(actual)
	suite.Error(actualErr)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func TestHandleOrganizationUpdateSuite(t *testing.T) {
	suite.Run(t, new(HandleOrganizationUpdateSuite))
}

func TestSendConfirmationNotificationSuite(t *testing.T) {
	suite.Run(t, new(SendConfirmationNotificationSuite))
}
