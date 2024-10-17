CREATE TABLE IF NOT EXISTS circulator_driver (
  driver_id uuid DEFAULT gen_random_uuid(),
  first_name VARCHAR(250) NOT NULL,
  middle_name VARCHAR(250) NULL,
  last_name VARCHAR(250) NOT NULL,
  suffix VARCHAR(250) NULL,
  created_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  last_modified_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  PRIMARY KEY (driver_id)
);