ALTER table chat.patients
RENAME COLUMN medicaid_id TO insurance_id;

ALTER index chat.idx_patient_medicaid_id RENAME TO idx_patient_insurance_id;

Alter table chat.patients
RENAME constraint patients_medicaid_id_key TO patients_insurance_id_key;