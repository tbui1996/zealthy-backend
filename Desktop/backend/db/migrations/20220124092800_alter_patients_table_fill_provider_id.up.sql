UPDATE chat.patients
SET provider_id = subquery.chat_user_id
FROM (
         SELECT DISTINCT chat.session_users.user_id       AS chat_user_id,
                         chat.session_users.session_id    AS chat_session_id,
                         chat.session_patients.patient_id AS chat_patient_id
         FROM chat.session_users
                  INNER JOIN chat.session_patients ON chat.session_users.session_id = chat.session_patients.session_id
         WHERE user_id NOT LIKE 'Okta_%'
     ) AS subquery
WHERE chat.patients.id = subquery.chat_patient_id;