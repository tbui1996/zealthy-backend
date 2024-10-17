package response

import (
	"time"
)

type PatientResponse struct {
	PatientId                  string     `json:"patientId"`
	InsuranceId                string     `json:"insuranceId"`
	FirstName                  string     `json:"patientFirstName"`
	MiddleName                 string     `json:"patientMiddleName"`
	LastName                   string     `json:"patientLastName"`
	Suffix                     string     `json:"patientSuffix"`
	DateOfBirth                *time.Time `json:"patientDateOfBirth"`
	PrimaryLanguage            string     `json:"patientPrimaryLanguage"`
	PreferredGender            string     `json:"patientPreferredGender"`
	EmailAddress               string     `json:"patientEmailAddress"`
	HomePhone                  string     `json:"patientHomePhone"`
	HomeLivingArrangement      string     `json:"patientHomeLivingArrangement"`
	HomeAddress1               string     `json:"patientHomeAddress1"`
	HomeAddress2               string     `json:"patientHomeAddress2"`
	HomeCity                   string     `json:"patientHomeCity"`
	HomeCounty                 string     `json:"patientHomeCounty"`
	HomeState                  string     `json:"patientHomeState"`
	HomeZip                    string     `json:"patientHomeZip"`
	SignedCirculoConsentForm   bool       `json:"patientSignedCirculoConsentForm"`
	CirculoConsentFormLink     string     `json:"patientCirculoConsentFormLink"`
	SignedStationMDConsentForm bool       `json:"patientSignedStationMDConsentForm"`
	StationMDConsentFormLink   string     `json:"patientStationMDConsentFormLink"`
	CompletedGoSheet           bool       `json:"patientCompletedGoSheet"`
	MarkedAsActive             bool       `json:"patientMarkedAsActive"`
	CreatedTimestamp           *time.Time `json:"patientCreatedTimestamp"`
	LastModifiedTimestamp      *time.Time `json:"patientLastModifiedTimestamp"`
}
