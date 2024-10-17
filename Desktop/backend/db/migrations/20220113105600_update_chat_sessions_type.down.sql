INSERT INTO chat.session_descriptors (session_id, name, value)
SELECT id, 'DESCRIPTION', chat_type from chat.sessions;