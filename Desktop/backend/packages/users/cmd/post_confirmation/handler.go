package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"go.uber.org/zap"
)

type HandlerInput struct {
	UserName     string
	PoolID       string
	Idp          cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Logger       *zap.Logger
	DefaultGroup string
}

func handler(input HandlerInput) *exception.SonarError {
	input.Logger.Info(fmt.Sprintf("Adding user to default group (%s)", input.UserName))
	_, err := input.Idp.AdminAddUserToGroup(&cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(input.DefaultGroup),
		Username:   aws.String(input.UserName),
		UserPoolId: aws.String(input.PoolID),
	})

	if err != nil {
		errMessage := fmt.Errorf("something went wrong adding user (%s) to group (%s): (%s)", input.UserName, input.DefaultGroup, err)
		input.Logger.Error(errMessage.Error())
		return exception.NewSonarError(http.StatusBadRequest, err.Error())
	}

	return nil
}
