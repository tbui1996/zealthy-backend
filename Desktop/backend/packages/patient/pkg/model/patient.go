package model

import (
	"time"
)

type Patient struct {
	PatientId                  string     `gorm:"primaryKey;column:patient_id"`
	InsuranceId                string     `gorm:"column:insurance_id"`
	FirstName                  string     `gorm:"column:first_name"`
	MiddleName                 string     `gorm:"column:middle_name"`
	LastName                   string     `gorm:"column:last_name"`
	Suffix                     string     `gorm:"column:suffix"`
	DateOfBirth                time.Time  `gorm:"column:date_of_birth"`
	PrimaryLanguage            string     `gorm:"column:primary_language"`
	PreferredGender            string     `gorm:"column:preferred_gender"`
	EmailAddress               string     `gorm:"column:email_address"`
	HomePhone                  string     `gorm:"column:home_phone"`
	HomeLivingArrangement      string     `gorm:"column:home_living_arrangement"`
	HomeAddress1               string     `gorm:"column:home_address_1"`
	HomeAddress2               string     `gorm:"column:home_address_2"`
	HomeCity                   string     `gorm:"column:home_city"`
	HomeCounty                 string     `gorm:"column:home_county"`
	HomeState                  string     `gorm:"column:home_state"`
	HomeZip                    string     `gorm:"column:home_zip"`
	SignedCirculoConsentForm   bool       `gorm:"column:signed_circulo_consent_form"`
	CirculoConsentFormLink     string     `gorm:"column:circulo_consent_form_link"`
	SignedStationMDConsentForm bool       `gorm:"column:signed_stationmd_consent_form"`
	StationMDConsentFormLink   string     `gorm:"column:stationmd_consent_form_link"`
	CompletedGoSheet           bool       `gorm:"column:completed_go_sheet"`
	MarkedAsActive             bool       `gorm:"column:marked_as_active"`
	CreatedTimestamp           *time.Time `gorm:"column:created_timestamp"`
	LastModifiedTimestamp      *time.Time `gorm:"column:last_modified_timestamp"`
}

func NewPatient(patient Patient) *Patient {
	return &Patient{
		InsuranceId:              patient.InsuranceId,
		FirstName:                patient.FirstName,
		MiddleName:               patient.MiddleName,
		LastName:                 patient.LastName,
		Suffix:                   patient.Suffix,
		DateOfBirth:              patient.DateOfBirth,
		PrimaryLanguage:          patient.PrimaryLanguage,
		PreferredGender:          patient.PreferredGender,
		EmailAddress:             patient.EmailAddress,
		HomePhone:                patient.HomePhone,
		HomeLivingArrangement:    patient.HomeLivingArrangement,
		HomeAddress1:             patient.HomeAddress1,
		HomeAddress2:             patient.HomeAddress2,
		HomeCity:                 patient.HomeCity,
		HomeCounty:               patient.HomeCounty,
		HomeState:                patient.HomeState,
		HomeZip:                  patient.HomeZip,
		CirculoConsentFormLink:   patient.CirculoConsentFormLink,
		StationMDConsentFormLink: patient.StationMDConsentFormLink,
	}
}

func (patient *Patient) IsNew() bool {
	return patient.CreatedTimestamp == nil
}
