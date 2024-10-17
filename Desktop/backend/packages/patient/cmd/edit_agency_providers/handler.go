package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AgencyProviderRepository
	logger *zap.Logger
}

func Handler(input request.EditAgencyProviderRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	agencyProvider, err := deps.repo.Find(input.AgencyProviderId)
	if err != nil {
		deps.logger.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
		return err.ToSonarApiGatewayResponse()
	}
	now := time.Now()
	agencyProvider.NationalProviderId = input.NationalProviderId
	agencyProvider.FirstName = input.FirstName
	agencyProvider.MiddleName = input.MiddleName
	agencyProvider.LastName = input.LastName
	agencyProvider.Suffix = input.Suffix
	agencyProvider.DoddNumber = input.DoddNumber
	agencyProvider.BusinessName = input.BusinessName
	agencyProvider.BusinessAddress1 = input.BusinessAddress1
	agencyProvider.BusinessAddress2 = input.BusinessAddress2
	agencyProvider.BusinessCity = input.BusinessCity
	agencyProvider.BusinessZip = input.BusinessZip
	agencyProvider.LastModifiedTimestamp = &now
	agencyProvider.BusinessTIN = input.BusinessTIN
	agencyProvider.BusinessState = input.BusinessState

	err = deps.repo.Save(agencyProvider)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
