package main

import (
	"net/http"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HandlerInput struct {
	UserName         string
	PoolID           string
	DB               *gorm.DB
	Logger           *zap.Logger
	OrganizationName string
}

type ExternalUser struct {
	ID                         string
	ExternalUserOrganizationId int
}

type ExternalUserOrganization struct {
	ID   int
	Name string
}

func handler(input HandlerInput) *exception.SonarError {

	organization := ExternalUserOrganization{Name: input.OrganizationName}
	res := input.DB.Where("name = ?", input.OrganizationName).FirstOrCreate(&organization)

	if res.Error != nil {
		return exception.NewSonarError(http.StatusInternalServerError, res.Error.Error())
	}

	input.Logger.Info("copying user to db")
	result := input.DB.Create(&ExternalUser{
		ID:                         input.UserName,
		ExternalUserOrganizationId: organization.ID,
	})

	if result.Error != nil {
		return exception.NewSonarError(http.StatusBadRequest, "could not insert user id in db: "+result.Error.Error())
	}

	return nil
}
