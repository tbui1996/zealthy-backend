CREATE TABLE IF NOT EXISTS form.form_sents (
  id SERIAL PRIMARY KEY,
  form_id INT NOT NULL,
  sent DATE NOT NULL,
  FOREIGN KEY (form_id)
      REFERENCES form.forms (id)
);