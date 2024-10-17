package request

type CreateAppointmentRequest struct {
	FirstName                     string
	LastName                      string
	PatientId                     string
	AgencyProviderId              string
	AppointmentScheduled          string
	AppointmentStatus             string
	AppointmentPurpose            string
	PatientChiefComplaint         string
	BusinessName                  string
	CirculatorDriverFullName      string
	PatientSystolicBloodPressure  int
	PatientDiastolicBloodPressure int
	PatientRespirationsPerMinute  int
	PatientPulseBeatsPerMinute    int
	PatientWeightLbs              int
	AppointmentOtherPurpose       string
	AppointmentNotes              string
}
