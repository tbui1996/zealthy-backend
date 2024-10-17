package request

type CreatePatientRequest struct {
	InsuranceId                     string
	PatientFirstName                string
	PatientMiddleName               string
	PatientLastName                 string
	PatientSuffix                   string
	PatientDateOfBirth              string
	PatientPrimaryLanguage          string
	PatientPreferredGender          string
	PatientEmailAddress             string
	PatientHomePhone                string
	PatientHomeLivingArrangement    string
	PatientHomeAddress1             string
	PatientHomeAddress2             string
	PatientHomeCity                 string
	PatientHomeCounty               string
	PatientHomeState                string
	PatientHomeZip                  string
	PatientCirculoConsentFormLink   string
	PatientStationMDConsentFormLink string
}
