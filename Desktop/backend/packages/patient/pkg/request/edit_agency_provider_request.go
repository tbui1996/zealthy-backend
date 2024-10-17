package request

type EditAgencyProviderRequest struct {
	DoddNumber         string
	AgencyProviderId   string
	NationalProviderId string
	FirstName          string
	MiddleName         string
	LastName           string
	Suffix             string
	BusinessName       string
	BusinessTIN        string
	BusinessAddress1   string
	BusinessAddress2   string
	BusinessCity       string
	BusinessState      string
	BusinessZip        string
}
