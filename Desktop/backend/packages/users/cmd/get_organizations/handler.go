package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/response"
	"go.uber.org/zap"
	"net/http"
)

type GetOrganizationsRequest struct {
	Logger   *zap.Logger
	Registry mapper.RegistryAPI
}

func handler(req GetOrganizationsRequest) ([]response.Organizations, error) {
	res, err := req.Registry.ExternalUserOrganization().FindAll()

	if err != nil {
		return nil, err
	}

	orgs := make([]response.Organizations, 0)
	for _, value := range res {
		if value == nil {
			req.Logger.Debug("got a nil organization in response from DB")
			continue
		}

		orgs = append(orgs, response.Organizations{ID: value.ID, Name: value.Name()})
	}

	return orgs, nil
}

func errorResponse(err error, message string, logger *zap.Logger) (events.APIGatewayV2HTTPResponse, error) {
	errMsg := fmt.Errorf("%s %s", message, err.Error())
	logger.Error(errMsg.Error())
	return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg.Error())
}
