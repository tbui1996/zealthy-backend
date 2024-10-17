package response

import (
	"time"
)

type AppointmentResponse struct {
	AppointmentId                 string     `json:"appointmentId"`
	PatientId                     string     `json:"patientId"`
	AgencyProviderId              string     `json:"agencyProviderId"`
	CirculatorDriverFullName      string     `json:"circulatorDriverFullName"`
	AppointmentCreated            *time.Time `json:"appointmentCreated"`
	AppointmentScheduled          *time.Time `json:"appointmentScheduled"`
	AppointmentStatus             string     `json:"appointmentStatus"`
	AppointmentStatusChangedOn    *time.Time `json:"appointmentStatusChangedOn"`
	AppointmentPurpose            string     `json:"appointmentPurpose"`
	AppointmentOtherPurpose       string     `json:"appointmentOtherPurpose"`
	AppointmentNotes              string     `json:"appointmentNotes"`
	PatientDiastolicBloodPressure int        `json:"patientDiastolicBloodPressure"`
	PatientSystolicBloodPressure  int        `json:"patientSystolicBloodPressure"`
	PatientRespirationsPerMinute  int        `json:"patientRespirationsPerMinute"`
	PatientPulseBeatsPerMinute    int        `json:"patientPulseBeatsPerMinute"`
	PatientWeightLbs              int        `json:"patientWeightLbs"`
	PatientChiefComplaint         string     `json:"patientChiefComplaint"`
	CreatedTimestamp              *time.Time `json:"createdTimestamp"`
	LastModifiedTimestamp         *time.Time `json:"lastModifiedTimestamp"`
	FirstName                     string     `json:"firstName"`
	MiddleName                    string     `json:"middleName"`
	LastName                      string     `json:"lastName"`
	ProviderFullName              string     `json:"providerFullName"`
	Suffix                        string     `json:"suffix"`
	DateOfBirth                   *time.Time `json:"dateOfBirth"`
	PrimaryLanguage               string     `json:"primaryLanguage"`
	PreferredGender               string     `json:"preferredGender"`
	EmailAddress                  string     `json:"emailAddress"`
	HomeAddress1                  string     `json:"homeAddress1"`
	HomeAddress2                  string     `json:"homeAddress2"`
	HomeCity                      string     `json:"homeCity"`
	HomeState                     string     `json:"homeState"`
	HomeZip                       string     `json:"homeZip"`
	SignedCirculoConsentForm      bool       `json:"signedCirculoConsentForm"`
	CirculoConsentFormLink        string     `json:"circuloConsentFormLink"`
	SignedStationMDConsentForm    bool       `json:"signedStationMDConsentForm"`
	StationMDConsentFormLink      string     `json:"stationMDConsentFormLink"`
	CompletedGoSheet              bool       `json:"completedGoSheet"`
	MarkedAsActive                bool       `json:"markedAsActive"`
	NationalProviderId            string     `json:"nationalProviderId"`
	BusinessName                  string     `json:"businessName"`
	BusinessTIN                   string     `json:"businessTIN"`
	BusinessAddress1              string     `json:"businessAddress1"`
	BusinessAddress2              string     `json:"businessAddress2"`
	BusinessCity                  string     `json:"businessCity"`
	BusinessState                 string     `json:"businessState"`
	BusinessZip                   string     `json:"businessZip"`
	HomePhone                     string     `json:"patientHomePhone"`
	HomeLivingArrangement         string     `json:"patientHomeLivingArrangement"`
	HomeCounty                    string     `json:"patientHomeCounty"`
	InsuranceId                   string     `json:"insuranceId"`
}
