package main

import (
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dto"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"go.uber.org/zap"
)

type HandlerInput struct {
	QueryStringParameters map[string]string
	Registry              mapper.RegistryAPI
	Logger                *zap.Logger
}

func handler(input HandlerInput) ([]*dto.ExternalUser, error) {
	mapper := input.Registry.ExternalUser()

	users, err := mapper.FindAll()

	if err != nil {
		return nil, exception.NewSonarError(http.StatusInternalServerError, "unable to find all users"+err.Error())
	}

	return dto.ExternalUsersFromModels(users), nil
}
