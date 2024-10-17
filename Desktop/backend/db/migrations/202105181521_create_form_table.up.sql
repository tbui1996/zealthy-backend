CREATE TABLE IF NOT EXISTS form.forms (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(255) NULL,
    created DATE NOT NULL
);