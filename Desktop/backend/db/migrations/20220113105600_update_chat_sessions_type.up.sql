UPDATE chat.sessions SET chat_type = subquery.value::text::chat_type
FROM (SELECT session_id, value FROM  chat.session_descriptors WHERE name = 'DESCRIPTION' ) AS subquery
WHERE chat.sessions.id = subquery.session_id;

DELETE FROM chat.session_descriptors WHERE name = 'DESCRIPTION';