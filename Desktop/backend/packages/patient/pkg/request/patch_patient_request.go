package request

type PatchPatientRequest struct {
	PatientId                         string
	InsuranceId                       string
	PatientFirstName                  string
	PatientMiddleName                 string
	PatientLastName                   string
	PatientSuffix                     string
	PatientDateOfBirth                string
	PatientPrimaryLanguage            string
	PatientPreferredGender            string
	PatientEmailAddress               string
	PatientHomePhone                  string
	PatientHomeLivingArrangement      string
	PatientHomeAddress1               string
	PatientHomeAddress2               string
	PatientHomeCity                   string
	PatientHomeCounty                 string
	PatientHomeState                  string
	PatientHomeZip                    string
	PatientSignedCirculoConsentForm   bool
	PatientCirculoConsentFormLink     string
	PatientSignedStationMDConsentForm bool
	PatientStationMDConsentFormLink   string
	PatientCompletedGoSheet           bool
	PatientMarkedAsActive             bool
}
