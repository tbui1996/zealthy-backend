ALTER TABLE form.forms
ADD COLUMN creator_id VARCHAR(25),
ADD COLUMN creator VARCHAR(100);

UPDATE form.forms
SET creator_id='UNKNOWN',
    creator='UNKNOWN';

ALTER TABLE form.forms
ALTER COLUMN creator_id SET NOT NULL,
ALTER COLUMN creator SET NOT NULL;