package response

import (
	"time"
)

type AgencyProviderResponse struct {
	AgencyProviderId      string     `json:"agencyProviderId"`
	DoddNumber            string     `json:"doddNumber"`
	NationalProviderId    string     `json:"nationalProviderId"`
	FirstName             string     `json:"firstName"`
	MiddleName            string     `json:"middleName"`
	LastName              string     `json:"lastName"`
	Suffix                string     `json:"suffix"`
	BusinessName          string     `json:"businessName"`
	BusinessTIN           string     `json:"businessTIN"`
	BusinessAddress1      string     `json:"businessAddress1"`
	BusinessAddress2      string     `json:"businessAddress2"`
	BusinessCity          string     `json:"businessCity"`
	BusinessState         string     `json:"businessState"`
	BusinessZip           string     `json:"businessZip"`
	CreatedTimestamp      *time.Time `json:"createdTimestamp"`
	LastModifiedTimestamp *time.Time `json:"lastModifiedTimestamp"`
}
