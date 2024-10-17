ALTER TABLE agency_provider
ADD column dodd_number VARCHAR(250) NOT NULL UNIQUE,
ADD CONSTRAINT national_provider_id_unique UNIQUE (national_provider_id);

ALTER TABLE agency_provider
ALTER COLUMN national_provider_id SET DEFAULT NULL;