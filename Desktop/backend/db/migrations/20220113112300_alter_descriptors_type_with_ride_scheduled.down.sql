AlTER TYPE SESSION_DESCRIPTOR RENAME TO SESSION_DESCRIPTOR_TEMP;
CREATE TYPE SESSION_DESCRIPTOR AS ENUM ('TOPIC', 'STARRED');
ALTER TABLE chat.session_descriptors ALTER COLUMN name TYPE SESSION_DESCRIPTOR using name::text::session_descriptor;
DROP TYPE SESSION_DESCRIPTOR_TEMP;
