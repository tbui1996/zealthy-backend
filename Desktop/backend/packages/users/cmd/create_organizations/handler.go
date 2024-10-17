package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"go.uber.org/zap"
)

type CreateOrganizationsRequest struct {
	Logger           *zap.Logger
	Registry         mapper.RegistryAPI
	OrganizationName string
}

func handler(req CreateOrganizationsRequest) (*model.ExternalUserOrganization, error) {
	newOrg, err := req.Registry.ExternalUserOrganization().Insert(&iface.ExternalUserOrganizationInsertInput{Name: req.OrganizationName})

	if err != nil {
		return nil, err
	}

	return newOrg, nil
}

func errorResponse(err error, message string, logger *zap.Logger) (events.APIGatewayV2HTTPResponse, error) {
	errMsg := fmt.Errorf("%s %s", message, err.Error())
	logger.Error(errMsg.Error())
	return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg.Error())
}
