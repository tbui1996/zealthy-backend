CREATE TABLE IF NOT EXISTS chat.session_patients (
    session_id INT PRIMARY KEY,
    patient_id INT NOT NULL,
    FOREIGN KEY (patient_id)
        REFERENCES chat.patients (id),
    FOREIGN KEY (session_id)
       REFERENCES chat.sessions (id)
);