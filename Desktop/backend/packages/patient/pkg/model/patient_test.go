package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PatientModelSuite struct {
	suite.Suite
}

func (s *PatientModelSuite) Test__IsNew() {
	model := &Patient{}
	now := time.Now()

	s.True(model.IsNew())

	model.CreatedTimestamp = &now

	s.False(model.IsNew())
}

func (s *PatientModelSuite) Test__NewPatient() {
	now := time.Now()
	model := NewPatient(Patient{
		InsuranceId:                "2",
		FirstName:                  "thomas",
		MiddleName:                 "",
		LastName:                   "bui",
		Suffix:                     "MR",
		DateOfBirth:                now,
		PrimaryLanguage:            "",
		PreferredGender:            "",
		EmailAddress:               "",
		HomePhone:                  "",
		HomeLivingArrangement:      "",
		HomeAddress1:               "yes",
		HomeAddress2:               "",
		HomeCity:                   "no",
		HomeCounty:                 "moco",
		HomeState:                  "dc",
		HomeZip:                    "00002",
		SignedCirculoConsentForm:   false,
		CirculoConsentFormLink:     "",
		SignedStationMDConsentForm: false,
		StationMDConsentFormLink:   "",
		CompletedGoSheet:           false,
		MarkedAsActive:             false,
		CreatedTimestamp:           &now,
		LastModifiedTimestamp:      &now,
	})

	s.Equal(model.InsuranceId, "2")
	s.Equal(model.FirstName, "thomas")
	s.Equal(model.LastName, "bui")
	s.Equal(model.Suffix, "MR")
	s.Equal(model.HomeAddress1, "yes")
	s.Equal(model.HomeCity, "no")
	s.Equal(model.HomeCounty, "moco")
	s.Equal(model.HomeState, "dc")
	s.Equal(model.HomeZip, "00002")

}

func TestPatientModel(t *testing.T) {
	suite.Run(t, new(PatientModelSuite))
}
