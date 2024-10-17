CREATE TABLE IF NOT EXISTS agency_provider (
  agency_provider_id uuid DEFAULT gen_random_uuid(),
  national_provider_id VARCHAR(250) NULL,
  first_name VARCHAR(250) NULL,
  middle_name VARCHAR(250) NULL,
  last_name VARCHAR(250) NULL,
  suffix VARCHAR(250) NULL,
  business_name varchar(250) NULL,
  business_tin VARCHAR(250) NULL,
  business_address_1 VARCHAR(250) NULL,
  business_address_2 VARCHAR(250) NULL,
  business_city VARCHAR(250) NULL,
  business_state VARCHAR(250) NULL,
  business_zip VARCHAR(250) NULL,
  created_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  last_modified_timestamp TIMESTAMP without time zone NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  PRIMARY KEY (agency_provider_id)
);