package main

import (
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
	"go.uber.org/zap"
)

type UpdateUserDependencies struct {
	Logger   *zap.Logger
	Registry mapper.RegistryAPI
	Db       dynamodbiface.DynamoDBAPI
	Api      apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

func handler(req request.UpdateUserRequest, deps UpdateUserDependencies) error {
	user, err := deps.Registry.ExternalUser().Find(req.ID)
	if err != nil {
		return err
	}
	user.SetFirstName(req.FirstName)
	user.SetLastName(req.LastName)

	err = handleOrganizationUpdate(req.OrganizationID, UpdateOrganizationDependencies{
		User:     user,
		Registry: deps.Registry,
	})

	if err != nil {
		return err
	}

	if req.Group != user.Group() {
		user.SetGroup(req.Group)
		user.SetEnabled(true)
	}

	_, err = deps.Registry.ExternalUser().Update(user)

	if err != nil {
		return err
	}

	err = sendConfirmedNotification(user.Email, SendNotificationDependencies{
		Logger: deps.Logger,
		Db:     deps.Db,
		Api:    deps.Api,
	})

	if err != nil {
		return err
	}

	return nil
}
