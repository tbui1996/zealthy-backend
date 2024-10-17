CREATE TABLE IF NOT EXISTS chat.patients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    medicaid_id VARCHAR(13) NOT NULL UNIQUE,
    birthday TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_patient_medicaid_id ON chat.patients (medicaid_id);