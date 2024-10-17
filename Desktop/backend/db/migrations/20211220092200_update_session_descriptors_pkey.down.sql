CREATE UNIQUE INDEX IF NOT EXISTS desc_id_name_temp_idx ON chat.session_descriptors (session_id);
ALTER TABLE chat.session_descriptors DROP CONSTRAINT session_descriptors_pkey, ADD CONSTRAINT session_descriptors_pkey PRIMARY KEY USING INDEX desc_id_name_temp_idx;