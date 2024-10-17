package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/idp/iface"
	"go.uber.org/zap"
)

type HandleRevokeUserInput struct {
	SonarIDP iface.SonarIdentityProvider
	Logger   *zap.Logger
	UserID   string
}

type PostConnectionInput struct {
	ConnectionID         *string
	Data                 []byte
	Logger               *zap.Logger
	ExternalWebsocketAPI apigatewaymanagementapiiface.ApiGatewayManagementApiAPI
}

type GetSonarWebsocketConnectionInput struct {
	UserID  string
	Logger  *zap.Logger
	UsersDB dynamo.Database
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required, but was not provided")
	}

	return nil
}

func disableUser(input *HandleRevokeUserInput) *exception.SonarError {
	_, err := input.SonarIDP.AdminDisableUser(input.UserID)
	if err != nil {
		errMessage := fmt.Sprintf("something went wrong disabling user: (%s)", err)
		input.Logger.Error(errMessage)
		return exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return nil
}

func signOutUser(input *HandleRevokeUserInput) *exception.SonarError {
	_, err := input.SonarIDP.AdminUserGlobalSignOut(input.UserID)
	if err != nil {
		errMessage := fmt.Sprintf("something went wrong signing out user: (%s)", err)
		input.Logger.Error(errMessage)
		return exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return nil
}

func listGroups(input *HandleRevokeUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, *exception.SonarError) {
	adminListGroupsForUserOuptut, errResponse := input.SonarIDP.AdminListGroupsForUser(input.UserID)
	if errResponse != nil {
		return nil, exception.NewSonarError(http.StatusBadRequest, "something went wrong getting groups for user: "+errResponse.Error())
	}

	return adminListGroupsForUserOuptut, nil
}

func removeUserFromGroup(input *HandleRevokeUserInput, groupName string) *exception.SonarError {
	_, err := input.SonarIDP.AdminRemoveUserFromGroup(input.UserID, groupName)
	if err != nil {
		errMessage := fmt.Sprintf("something went wrong removing user from group: (%s)", err.Error())
		input.Logger.Error(errMessage)
		return exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return nil
}

func getUser(input *HandleRevokeUserInput) (*cognitoidentityprovider.AdminGetUserOutput, *exception.SonarError) {
	user, err := input.SonarIDP.AdminGetUser(input.UserID)
	if err != nil {
		errMessage := fmt.Sprintf("something went wrong getting user: (%s)", err)
		input.Logger.Error(errMessage)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return user, nil
}

func getSonarWebsocketConnection(input *GetSonarWebsocketConnectionInput) (*dynamodb.QueryOutput, *exception.SonarError) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(dynamo.SonarWebsocketConnections),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				S: aws.String(input.UserID),
			},
		},
		ProjectionExpression: aws.String("ConnectionId"),
	}

	results, err := input.UsersDB.Query(params)

	input.Logger.Info(fmt.Sprintf("Query API results: (%s)", results.Items))

	if err != nil {
		errMessage := fmt.Sprintf("Query API call failed: (%s)", err)
		input.Logger.Error(errMessage)
		return nil, exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return results, nil
}

func postConnection(input *PostConnectionInput) *exception.SonarError {
	connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: input.ConnectionID,
		Data:         input.Data,
	}

	_, err := input.ExternalWebsocketAPI.PostToConnection(connectionInput)
	if err != nil {
		errMessage := fmt.Sprintf("Error sending message: (%s)", err)
		input.Logger.Error(errMessage)
		return exception.NewSonarError(http.StatusInternalServerError, errMessage)
	}

	return nil
}
