package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/mocks"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type UpdateUserSuite struct {
	suite.Suite
	dependencies                   UpdateUserDependencies
	ExternalUser                   *iface.MockExternalUser
	externalUserOrganizationMapper *iface.MockExternalUserOrganization
	registry                       *mapper.MockRegistryAPI
	userToUpdate                   *model.ExternalUser
	organization                   *model.ExternalUserOrganization
	reqWithOrganization            request.UpdateUserRequest
	reqWithNoOrganization          request.UpdateUserRequest
	dynamoQueryInput               *dynamodb.QueryInput
	api                            apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

func (suite *UpdateUserSuite) SetupTest() {
	suite.registry = new(mapper.MockRegistryAPI)
	suite.ExternalUser = new(iface.MockExternalUser)
	suite.api = new(mocks.ApiGatewayManagementApiAPI)
	suite.externalUserOrganizationMapper = new(iface.MockExternalUserOrganization)
	suite.reqWithOrganization = request.UpdateUserRequest{
		ID:             "1",
		FirstName:      "Name",
		LastName:       "Last name",
		Group:          "external_supervisors",
		OrganizationID: 1,
	}

	suite.reqWithNoOrganization = request.UpdateUserRequest{
		ID:             "1",
		FirstName:      "Name",
		LastName:       "Last name",
		Group:          "external_supervisors",
		OrganizationID: 0,
	}

	suite.dependencies = UpdateUserDependencies{
		Registry: suite.registry,
		Logger:   zaptest.NewLogger(suite.T()),
		Db:       new(mocks.DynamoDBAPI),
		Api:      suite.api,
	}

	suite.userToUpdate = &model.ExternalUser{
		ID:       "1",
		Username: "Username",
		Email:    "name@circulohealth.com",
		Status:   "",
	}

	suite.organization = &model.ExternalUserOrganization{
		ID: 1,
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

func (suite *UpdateUserSuite) TestUpdateUserNameAndGroup_Success() {
	suite.registry.On("ExternalUser").Return(suite.ExternalUser)

	suite.ExternalUser.On("Find", "1").Return(suite.userToUpdate, nil)
	suite.ExternalUser.On("Update", suite.userToUpdate).Return(suite.userToUpdate, nil)

	suite.dependencies.Db.(*mocks.DynamoDBAPI).On("Query", suite.dynamoQueryInput).Return(&dynamodb.QueryOutput{}, nil)

	actualErr := handler(suite.reqWithNoOrganization, suite.dependencies)

	suite.NoError(actualErr)
	suite.ExternalUser.AssertCalled(suite.T(), "Find", "1")
	suite.ExternalUser.AssertCalled(suite.T(), "Update", suite.userToUpdate)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func (suite *UpdateUserSuite) TestUpdateUserNameAndGroup_FindFail() {
	suite.registry.On("ExternalUser").Return(suite.ExternalUser)

	suite.ExternalUser.On("Find", "1").Return(nil, errors.New("FAKE ERROR"))

	actualErr := handler(suite.reqWithNoOrganization, suite.dependencies)

	suite.Error(actualErr)
	suite.ExternalUser.AssertCalled(suite.T(), "Find", "1")
	suite.ExternalUser.AssertNotCalled(suite.T(), "Update", mock.Anything)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertNotCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func (suite *UpdateUserSuite) TestUpdateUserNameAndGroup_UpdateFail() {
	suite.registry.On("ExternalUser").Return(suite.ExternalUser)

	suite.ExternalUser.On("Find", "1").Return(suite.userToUpdate, nil)
	suite.ExternalUser.On("Update", suite.userToUpdate).Return(nil, errors.New("FAKE ERROR"))

	actualErr := handler(suite.reqWithNoOrganization, suite.dependencies)

	suite.Error(actualErr)
	suite.ExternalUser.AssertCalled(suite.T(), "Find", "1")
	suite.ExternalUser.AssertCalled(suite.T(), "Update", suite.userToUpdate)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertNotCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func (suite *UpdateUserSuite) TestUpdateOrganization_Success() {
	suite.registry.On("ExternalUser").Return(suite.ExternalUser)
	suite.ExternalUser.On("Find", "1").Return(suite.userToUpdate, nil)

	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Find", 1).Return(suite.organization, nil)

	suite.ExternalUser.On("Update", suite.userToUpdate).Return(suite.userToUpdate, nil)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).On("Query", suite.dynamoQueryInput).Return(&dynamodb.QueryOutput{}, nil)

	actualErr := handler(suite.reqWithOrganization, suite.dependencies)

	suite.NoError(actualErr)
	suite.ExternalUser.AssertCalled(suite.T(), "Find", "1")
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Find", 1)
	suite.ExternalUser.AssertCalled(suite.T(), "Update", suite.userToUpdate)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func (suite *UpdateUserSuite) TestUpdateOrganization_Fail() {
	suite.registry.On("ExternalUser").Return(suite.ExternalUser)
	suite.ExternalUser.On("Find", "1").Return(suite.userToUpdate, nil)

	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationMapper)
	suite.externalUserOrganizationMapper.On("Find", 1).Return(nil, errors.New("FAKE ERROR"))

	actualErr := handler(suite.reqWithOrganization, suite.dependencies)

	suite.Error(actualErr)
	suite.ExternalUser.AssertCalled(suite.T(), "Find", "1")
	suite.externalUserOrganizationMapper.AssertCalled(suite.T(), "Find", 1)
	suite.ExternalUser.AssertNotCalled(suite.T(), "Update", mock.Anything)
	suite.dependencies.Db.(*mocks.DynamoDBAPI).AssertNotCalled(suite.T(), "Query", suite.dynamoQueryInput)
}

func TestUpdateUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserSuite))
}
