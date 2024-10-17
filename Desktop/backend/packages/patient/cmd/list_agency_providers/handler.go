package main

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/response"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AgencyProviderRepository
	logger *zap.Logger
}

func Handler(deps HandlerDeps) (*[]response.AgencyProviderResponse, error) {
	agencyProviders, err := deps.repo.FindAll()

	if err != nil {
		deps.logger.Error(err.Error())
		return nil, fmt.Errorf("unable to find records in database")
	}

	responses := make([]response.AgencyProviderResponse, len(*agencyProviders))
	for i := range *agencyProviders {
		agencyProvider := (*agencyProviders)[i]
		responses[i] = response.AgencyProviderResponse{
			AgencyProviderId:      agencyProvider.AgencyProviderId,
			NationalProviderId:    agencyProvider.NationalProviderId,
			DoddNumber:            agencyProvider.DoddNumber,
			FirstName:             agencyProvider.FirstName,
			MiddleName:            agencyProvider.MiddleName,
			LastName:              agencyProvider.LastName,
			Suffix:                agencyProvider.Suffix,
			BusinessName:          agencyProvider.BusinessName,
			BusinessTIN:           agencyProvider.BusinessTIN,
			BusinessAddress1:      agencyProvider.BusinessAddress1,
			BusinessAddress2:      agencyProvider.BusinessAddress2,
			BusinessCity:          agencyProvider.BusinessCity,
			BusinessState:         agencyProvider.BusinessState,
			BusinessZip:           agencyProvider.BusinessZip,
			CreatedTimestamp:      agencyProvider.CreatedTimestamp,
			LastModifiedTimestamp: agencyProvider.LastModifiedTimestamp,
		}
	}

	return &responses, nil
}
