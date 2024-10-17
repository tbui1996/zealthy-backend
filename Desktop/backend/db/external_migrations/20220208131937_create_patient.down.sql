ALTER table patient 
DROP COLUMN insurance_id,
ADD column medicaid_id VARCHAR(250) NULL;