INSERT INTO chat.session_descriptors (session_id, name, value)
SELECT session_id, 'RIDE_SCHEDULED', ride_scheduled::TEXT FROM chat.session_statuses
WHERE ride_scheduled IS NOT NULL;