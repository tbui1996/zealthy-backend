CREATE TYPE CHAT_STATUS AS ENUM ('OPEN', 'CLOSED', 'PENDING');

CREATE TABLE IF NOT EXISTS chat.session_statuses (
   session_id INT PRIMARY KEY,
   status CHAT_STATUS NOT NULL,
   FOREIGN KEY (session_id)
       REFERENCES chat.sessions (id)
);
