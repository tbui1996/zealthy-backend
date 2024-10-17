package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/idp/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/request"
	"go.uber.org/zap"
)

type HandlerInput struct {
	Request  *request.RevokeAccessRequest
	Logger   *zap.Logger
	SonarIDP iface.SonarIdentityProvider
}

func Handler(input *HandlerInput) *exception.SonarError {
	input.Logger.Info(fmt.Sprintf("Requesting to remove user from groups, sign user out, and disable user: (%s) ", input.Request.Username))
	user, getUserErr := getUser(&HandleRevokeUserInput{
		SonarIDP: input.SonarIDP,
		Logger:   input.Logger,
		UserID:   input.Request.Username,
	})

	if getUserErr != nil {
		return getUserErr
	}
	if user == nil {
		return exception.NewSonarError(http.StatusBadRequest, "expected user to not be null: "+getUserErr.Error())
	}

	listGroups, listGroupsErr := listGroups(&HandleRevokeUserInput{
		SonarIDP: input.SonarIDP,
		Logger:   input.Logger,
		UserID:   input.Request.Username,
	})

	if listGroupsErr != nil {
		return listGroupsErr
	}
	if len(listGroups.Groups) == 1 {
		removeUserFromGroupErr := removeUserFromGroup(&HandleRevokeUserInput{
			SonarIDP: input.SonarIDP,
			Logger:   input.Logger,
			UserID:   input.Request.Username,
		}, *listGroups.Groups[0].GroupName)

		if removeUserFromGroupErr != nil {
			return removeUserFromGroupErr
		}
	} else {
		return exception.NewSonarError(http.StatusBadRequest, "expected list of groups per user to only have one group.")
	}

	signOutUserErr := signOutUser(&HandleRevokeUserInput{
		SonarIDP: input.SonarIDP,
		Logger:   input.Logger,
		UserID:   input.Request.Username,
	})

	if signOutUserErr != nil {
		return signOutUserErr
	}

	if *user.Enabled {
		disableUserErr := disableUser(&HandleRevokeUserInput{
			SonarIDP: input.SonarIDP,
			Logger:   input.Logger,
			UserID:   input.Request.Username,
		})

		if disableUserErr != nil {
			return disableUserErr
		}
	}
	return nil
}

type SendRevokedNotificationInput struct {
	UserID               string
	Logger               *zap.Logger
	UsersDB              dynamo.Database
	ExternalWebsocketAPI apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

func sendRevokedNotification(input *SendRevokedNotificationInput) *exception.SonarError {
	results, resultErr := getSonarWebsocketConnection(&GetSonarWebsocketConnectionInput{
		Logger:  input.Logger,
		UserID:  input.UserID,
		UsersDB: input.UsersDB,
	})

	if resultErr != nil {
		return resultErr
	}

	for _, item := range results.Items {
		connectionID := item["ConnectionId"].S
		itemLogger := input.Logger.With(zap.String("connectionID", *connectionID))
		itemLogger.Info(fmt.Sprintf("Sending message to: (%s)", *connectionID))
		connectionErr := postConnection(&PostConnectionInput{
			ConnectionID:         connectionID,
			Data:                 []byte("REVOKED"),
			Logger:               itemLogger,
			ExternalWebsocketAPI: input.ExternalWebsocketAPI,
		})

		if connectionErr != nil {
			return connectionErr
		}
	}

	return nil
}
