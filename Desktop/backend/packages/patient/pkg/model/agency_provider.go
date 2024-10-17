package model

import (
	"time"
)

type AgencyProvider struct {
	AgencyProviderId      string     `gorm:"primaryKey;column:agency_provider_id"`
	DoddNumber            string     `gorm:"column:dodd_number"`
	NationalProviderId    string     `gorm:"column:national_provider_id"`
	FirstName             string     `gorm:"column:first_name"`
	MiddleName            string     `gorm:"column:middle_name"`
	LastName              string     `gorm:"column:last_name"`
	Suffix                string     `gorm:"column:suffix"`
	BusinessName          string     `gorm:"column:business_name"`
	BusinessTIN           string     `gorm:"column:business_tin"`
	BusinessAddress1      string     `gorm:"column:business_address_1"`
	BusinessAddress2      string     `gorm:"column:business_address_2"`
	BusinessCity          string     `gorm:"column:business_city"`
	BusinessState         string     `gorm:"column:business_state"`
	BusinessZip           string     `gorm:"column:business_zip"`
	CreatedTimestamp      *time.Time `gorm:"column:created_timestamp"`
	LastModifiedTimestamp *time.Time `gorm:"column:last_modified_timestamp"`
}

func NewAgencyProvider(agencyProvider AgencyProvider) *AgencyProvider {
	return &AgencyProvider{
		AgencyProviderId:      agencyProvider.AgencyProviderId,
		DoddNumber:            agencyProvider.DoddNumber,
		NationalProviderId:    agencyProvider.NationalProviderId,
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

func (agencyProvider *AgencyProvider) IsNew() bool {
	return agencyProvider.CreatedTimestamp == nil
}
