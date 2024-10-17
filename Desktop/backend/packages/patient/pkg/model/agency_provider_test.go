package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AgencyProviderModelSuite struct {
	suite.Suite
}

func (s *AgencyProviderModelSuite) Test__IsNew() {
	model := &AgencyProvider{}
	now := time.Now()

	s.True(model.IsNew())

	model.CreatedTimestamp = &now

	s.False(model.IsNew())
}

func (s *AgencyProviderModelSuite) Test__NewAgencyProvider() {
	now := time.Now()
	model := NewAgencyProvider(AgencyProvider{
		AgencyProviderId:      "1",
		NationalProviderId:    "12",
		DoddNumber:            "13",
		FirstName:             "thomas",
		MiddleName:            "hehe",
		LastName:              "bui",
		Suffix:                "MR",
		BusinessName:          "circulo",
		BusinessTIN:           "7",
		BusinessAddress1:      "now",
		BusinessAddress2:      "",
		BusinessCity:          "columbus",
		BusinessState:         "oh",
		BusinessZip:           "00001",
		CreatedTimestamp:      &now,
		LastModifiedTimestamp: &now,
	})

	s.Equal(model.AgencyProviderId, "1")
	s.Equal(model.NationalProviderId, "12")
	s.Equal(model.DoddNumber, "13")
	s.Equal(model.FirstName, "thomas")
	s.Equal(model.MiddleName, "hehe")
	s.Equal(model.LastName, "bui")
	s.Equal(model.Suffix, "MR")
	s.Equal(model.BusinessName, "circulo")
	s.Equal(model.BusinessTIN, "7")
	s.Equal(model.BusinessAddress1, "now")
	s.Equal(model.BusinessAddress2, "")
	s.Equal(model.BusinessCity, "columbus")
	s.Equal(model.BusinessZip, "00001")
	s.Equal(model.BusinessState, "oh")
	s.Equal(model.CreatedTimestamp, &now)
	s.Equal(model.LastModifiedTimestamp, &now)

}

func TestAgencyProviderModel(t *testing.T) {
	suite.Run(t, new(AgencyProviderModelSuite))
}
