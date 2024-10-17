ALTER index chat.idx_patient_insurance_id RENAME TO idx_patient_medicaid_id;
Alter table chat.patients
RENAME constraint patients_insurance_id_key TO patients_medicaid_id_key;

ALTER table chat.patients
RENAME COLUMN insurance_id TO medicaid_id;