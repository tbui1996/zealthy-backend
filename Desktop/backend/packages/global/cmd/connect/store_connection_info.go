package main

import (
	"github.com/circulohealth/sonar-backend/packages/global/pkg/model"
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
)

type StoreConnectionInfoInput struct {
	Context validate.ConnectionContext
	DB      dynamo.Database
}

func storeConnectionInfo(input StoreConnectionInfoInput) *exception.SonarError {
	item := model.ConnectionItem{
		ConnectionId: input.Context.ConnectionID,
		UserID:       input.Context.UserID,
		CognitoGroup: input.Context.CognitoGroup,
	}

	err := input.DB.Create(item)
	if err != nil {
		return exception.NewSonarError(http.StatusBadRequest, "calling PutItem for connection item: "+err.Error())
	}

	return nil
}
