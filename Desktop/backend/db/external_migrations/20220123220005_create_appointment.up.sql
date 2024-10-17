CREATE TABLE IF NOT EXISTS appointments (
  appointment_id uuid DEFAULT gen_random_uuid(),
  patient_id uuid NOT NULL,
  agency_provider_id uuid NOT NULL,
  circulator_driver_fullname VARCHAR(250) NOT NULL,
  appointment_created timestamp without time zone NOT NULL,
  appointment_scheduled timestamp without time zone NOT NULL,
  appointment_status varchar(250) NOT NULL,
  appointment_status_changed_on timestamp without time zone NOT NULL,
  appointment_purpose varchar(250) NOT NULL,
  appointment_other_purpose varchar(500),
  appointment_notes varchar(500),
  patient_diastolic_blood_pressure integer NOT NULL DEFAULT -1,
  patient_systolic_blood_pressure integer NOT NULL DEFAULT -1,
  patient_respirations_per_minute integer NOT NULL DEFAULT -1,
  patient_pulse_beats_per_minute integer NOT NULL DEFAULT -1,
  patient_weight_lbs integer NOT NULL DEFAULT -1,
  patient_chief_complaint varchar(500) NULL,
  created_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  last_modified_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  PRIMARY KEY (appointment_id),
  FOREIGN KEY(patient_id)
    REFERENCES patient (patient_id),
  FOREIGN KEY(agency_provider_id)
    REFERENCES agency_provider (agency_provider_id)
);
