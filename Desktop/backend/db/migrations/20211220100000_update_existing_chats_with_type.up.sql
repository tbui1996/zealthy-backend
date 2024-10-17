INSERT INTO chat.session_descriptors (session_id, name, value)
SELECT id, 'DESCRIPTION', 'CIRCULATOR' from chat.sessions;