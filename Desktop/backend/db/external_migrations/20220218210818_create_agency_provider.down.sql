ALTER TABLE agency_provider
DROP COLUMN IF EXISTS dodd_number,
DROP CONSTRAINT national_provider_id_unique;