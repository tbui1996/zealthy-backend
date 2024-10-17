package model

import (
	"time"
)

type JoinResult struct {
	AppointmentId                 string     `gorm:"primaryKey;column:appointment_id"`
	PatientId                     string     `gorm:"column:patient_id"`
	AgencyProviderId              string     `gorm:"column:agency_provider_id"`
	CirculatorDriverFullName      string     `gorm:"column:circulator_driver_fullname"`
	AppointmentCreated            *time.Time `gorm:"column:appointment_created"`
	AppointmentScheduled          *time.Time `gorm:"column:appointment_scheduled"`
	AppointmentStatus             string     `gorm:"column:appointment_status"`
	AppointmentStatusChangedOn    *time.Time `gorm:"column:appointment_status_changed_on"`
	AppointmentPurpose            string     `gorm:"column:appointment_purpose"`
	AppointmentOtherPurpose       string     `gorm:"column:appointment_other_purpose"`
	AppointmentNotes              string     `gorm:"column:appointment_notes"`
	PatientDiastolicBloodPressure int        `gorm:"column:patient_diastolic_blood_pressure"`
	PatientSystolicBloodPressure  int        `gorm:"column:patient_systolic_blood_pressure"`
	PatientRespirationsPerMinute  int        `gorm:"column:patient_respirations_per_minute"`
	PatientPulseBeatsPerMinute    int        `gorm:"column:patient_pulse_beats_per_minute"`
	PatientWeightLbs              int        `gorm:"column:patient_weight_lbs"`
	PatientChiefComplaint         string     `gorm:"patient_chief_complaint"`
	CreatedTimestamp              *time.Time `gorm:"created_timestamp"`
	LastModifiedTimestamp         *time.Time `gorm:"last_modified_timestamp"`
	FirstName                     string     `gorm:"column:first_name"`
	MiddleName                    string     `gorm:"column:middle_name"`
	LastName                      string     `gorm:"column:last_name"`
	ProviderFullName              string     `gorm:"column:provider_fullname"`
	Suffix                        string     `gorm:"column:suffix"`
	DateOfBirth                   time.Time  `gorm:"column:date_of_birth"`
	PrimaryLanguage               string     `gorm:"column:primary_language"`
	PreferredGender               string     `gorm:"column:preferred_gender"`
	EmailAddress                  string     `gorm:"column:email_address"`
	HomeAddress1                  string     `gorm:"column:home_address_1"`
	HomeAddress2                  string     `gorm:"column:home_address_2"`
	HomeCity                      string     `gorm:"column:home_city"`
	HomeState                     string     `gorm:"column:home_state"`
	HomeZip                       string     `gorm:"column:home_zip"`
	SignedCirculoConsentForm      bool       `gorm:"column:signed_circulo_consent_form"`
	CirculoConsentFormLink        string     `gorm:"column:circulo_consent_form_link"`
	SignedStationMDConsentForm    bool       `gorm:"column:signed_stationmd_consent_form"`
	StationMDConsentFormLink      string     `gorm:"column:stationmd_consent_form_link"`
	CompletedGoSheet              bool       `gorm:"column:completed_go_sheet"`
	MarkedAsActive                bool       `gorm:"column:marked_as_active"`
	NationalProviderId            string     `gorm:"column:national_provider_id"`
	BusinessName                  string     `gorm:"column:business_name"`
	BusinessTIN                   string     `gorm:"column:business_tin"`
	BusinessAddress1              string     `gorm:"column:business_address_1"`
	BusinessAddress2              string     `gorm:"column:business_address_2"`
	BusinessCity                  string     `gorm:"column:business_city"`
	BusinessState                 string     `gorm:"column:business_state"`
	BusinessZip                   string     `gorm:"column:business_zip"`
	HomePhone                     string     `gorm:"column:home_phone"`
	HomeLivingArrangement         string     `gorm:"column:home_living_arrangement"`
	HomeCounty                    string     `gorm:"column:home_county"`
	InsuranceId                   string     `gorm:"column:insurance_id"`
}

type Appointment struct {
	AppointmentId                 string     `gorm:"primaryKey;column:appointment_id"`
	PatientId                     string     `gorm:"column:patient_id"`
	AgencyProviderId              string     `gorm:"column:agency_provider_id"`
	CirculatorDriverFullName      string     `gorm:"column:circulator_driver_fullname"`
	AppointmentCreated            *time.Time `gorm:"column:appointment_created"`
	AppointmentScheduled          string     `gorm:"column:appointment_scheduled"`
	AppointmentStatus             string     `gorm:"column:appointment_status"`
	AppointmentStatusChangedOn    *time.Time `gorm:"column:appointment_status_changed_on"`
	AppointmentPurpose            string     `gorm:"column:appointment_purpose"`
	AppointmentOtherPurpose       string     `gorm:"column:appointment_other_purpose"`
	AppointmentNotes              string     `gorm:"column:appointment_notes"`
	PatientDiastolicBloodPressure int        `gorm:"column:patient_diastolic_blood_pressure"`
	PatientSystolicBloodPressure  int        `gorm:"column:patient_systolic_blood_pressure"`
	PatientRespirationsPerMinute  int        `gorm:"column:patient_respirations_per_minute"`
	PatientPulseBeatsPerMinute    int        `gorm:"column:patient_pulse_beats_per_minute"`
	PatientWeightLbs              int        `gorm:"column:patient_weight_lbs"`
	PatientChiefComplaint         string     `gorm:"patient_chief_complaint"`
	CreatedTimestamp              *time.Time `gorm:"created_timestamp"`
	LastModifiedTimestamp         *time.Time `gorm:"last_modified_timestamp"`
}

func NewAppointment(appointment Appointment) *Appointment {
	return &Appointment{
		AppointmentId:                 appointment.AppointmentId,
		PatientId:                     appointment.PatientId,
		AgencyProviderId:              appointment.AgencyProviderId,
		CirculatorDriverFullName:      appointment.CirculatorDriverFullName,
		AppointmentCreated:            appointment.AppointmentCreated,
		AppointmentScheduled:          appointment.AppointmentScheduled,
		AppointmentStatus:             appointment.AppointmentStatus,
		AppointmentStatusChangedOn:    appointment.AppointmentStatusChangedOn,
		AppointmentPurpose:            appointment.AppointmentPurpose,
		AppointmentOtherPurpose:       appointment.AppointmentOtherPurpose,
		AppointmentNotes:              appointment.AppointmentNotes,
		PatientDiastolicBloodPressure: appointment.PatientDiastolicBloodPressure,
		PatientSystolicBloodPressure:  appointment.PatientSystolicBloodPressure,
		PatientRespirationsPerMinute:  appointment.PatientRespirationsPerMinute,
		PatientPulseBeatsPerMinute:    appointment.PatientPulseBeatsPerMinute,
		PatientWeightLbs:              appointment.PatientWeightLbs,
		PatientChiefComplaint:         appointment.PatientChiefComplaint,
		CreatedTimestamp:              appointment.CreatedTimestamp,
		LastModifiedTimestamp:         appointment.LastModifiedTimestamp,
	}
}

func (appointment *Appointment) IsNew() bool {
	return appointment.CreatedTimestamp == nil
}
