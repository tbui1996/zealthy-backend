package main

import (
	"errors"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/model"
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/common/dynamo"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/global/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func inputStoreConnectionInfoBuilder(userID, connectionID, group string, dbResponse interface{}) StoreConnectionInfoInput {
	mockDB := new(dynamo.MockDatabase)

	item := model.ConnectionItem{
		ConnectionId: connectionID,
		UserID:       userID,
		CognitoGroup: group,
	}

	mockDB.On("Create", item).Return(dbResponse)

	return StoreConnectionInfoInput{
		Context: validate.ConnectionContext{
			ConnectionID: connectionID,
			UserID:       userID,
			CognitoGroup: group,
		},
		DB: mockDB,
	}
}

func TestStoreConnectionInfo(t *testing.T) {
	tests := []struct {
		input       StoreConnectionInfoInput
		expectedErr *exception.SonarError
	}{
		{
			// no user id
			input:       inputStoreConnectionInfoBuilder("", "connectionid", "", errors.New("no UserID")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling PutItem for connection item: no UserID"),
		},
		{
			// no connection id
			input:       inputStoreConnectionInfoBuilder("userid", "", "", errors.New("no ConnectionID")),
			expectedErr: exception.NewSonarError(http.StatusBadRequest, "calling PutItem for connection item: no ConnectionID"),
		},
		{
			// valid store, no error
			input:       inputStoreConnectionInfoBuilder("userid", "connectionid", "group", nil),
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		actualErr := storeConnectionInfo(test.input)
		assert.Equal(t, test.expectedErr, actualErr)

		assert.IsType(t, new(dynamo.MockDatabase), test.input.DB)

		mockDB, _ := test.input.DB.(*dynamo.MockDatabase)

		mockDB.AssertCalled(t, "Create", model.ConnectionItem{
			ConnectionId: test.input.Context.ConnectionID,
			UserID:       test.input.Context.UserID,
			CognitoGroup: test.input.Context.CognitoGroup,
		})
	}
}
