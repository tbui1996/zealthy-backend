ALTER table patient 
DROP COLUMN medicaid_id,
ADD column insurance_id VARCHAR(250) NOT NULL UNIQUE;