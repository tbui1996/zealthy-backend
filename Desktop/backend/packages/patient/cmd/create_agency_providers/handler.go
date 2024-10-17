package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/data/iface"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/request"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	repo   iface.AgencyProviderRepository
	logger *zap.Logger
}

func Handler(input request.CreateAgencyProviderRequest, deps HandlerDeps) (events.APIGatewayV2HTTPResponse, error) {
	errMsg := "First/Last Name, DoDD Number, Business Address, Business Zip, and Business Name are all required"
	if input.FirstName == "" || input.LastName == "" || input.BusinessAddress1 == "" || input.BusinessZip == "" || input.BusinessName == "" || input.DoddNumber == "" {
		deps.logger.Error(errMsg)
		return exception.ErrorMessageApiGatewayV2(http.StatusBadRequest, errMsg)
	}

	agencyProvider := &model.AgencyProvider{
		NationalProviderId: input.NationalProviderId,
		FirstName:          input.FirstName,
		DoddNumber:         input.DoddNumber,
		MiddleName:         input.MiddleName,
		LastName:           input.LastName,
		Suffix:             input.Suffix,
		BusinessName:       input.BusinessName,
		BusinessTIN:        input.BusinessTIN,
		BusinessAddress1:   input.BusinessAddress1,
		BusinessAddress2:   input.BusinessAddress2,
		BusinessCity:       input.BusinessCity,
		BusinessState:      input.BusinessState,
		BusinessZip:        input.BusinessZip,
	}

	err := deps.repo.Save(agencyProvider)

	if err != nil {
		return err.ToSonarApiGatewayResponse()
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
