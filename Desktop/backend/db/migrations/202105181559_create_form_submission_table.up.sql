CREATE TABLE IF NOT EXISTS form.form_submissions (
  id SERIAL PRIMARY KEY,
  form_sent_id INT NOT NULL,
  FOREIGN KEY (form_sent_id)
      REFERENCES form.form_sents (id)
);