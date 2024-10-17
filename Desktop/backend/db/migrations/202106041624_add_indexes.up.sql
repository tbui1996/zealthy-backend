CREATE INDEX idx_inputs_form_id ON form.inputs (form_id);

CREATE INDEX idx_form_sents_form_id ON form.form_sents (form_id);

CREATE INDEX idx_form_submissions_form_sent_id ON form.form_submissions (form_sent_id);

CREATE INDEX idx_input_submissions_form_submission_id ON form.input_submissions (form_submission_id);

CREATE INDEX idx_form_discards_form_sent_id ON form.form_discards (form_sent_id);