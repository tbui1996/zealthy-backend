CREATE TABLE IF NOT EXISTS users.external_user_organizations (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE
);