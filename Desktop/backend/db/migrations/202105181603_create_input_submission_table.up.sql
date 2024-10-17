CREATE TABLE IF NOT EXISTS form.input_submissions (
    id SERIAL PRIMARY KEY,
    form_submission_id INT NOT NULL,
    input_id INT NOT NULL,
    response VARCHAR(255) NULL,
    FOREIGN KEY (form_submission_id)
        REFERENCES form.form_submissions (id),
    FOREIGN KEY (input_id)
        REFERENCES form.inputs (id)
);
