CREATE TABLE IF NOT EXISTS form.inputs (
  id SERIAL PRIMARY KEY,
  "order" INT DEFAULT 0,
  type VARCHAR(100) NOT NULL,
  label VARCHAR(255) NOT NULL,
  options TEXT [],
  form_id INT NOT NULL,
  FOREIGN KEY (form_id)
      REFERENCES form.forms (id)
);