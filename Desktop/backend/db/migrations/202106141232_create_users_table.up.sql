CREATE TABLE router.users (
  id SERIAL PRIMARY KEY,
  email VARCHAR NOT NULL UNIQUE,
  created_timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX email_idx ON router.users (email);