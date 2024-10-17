CREATE TABLE IF NOT EXISTS chat.session_last_messages (
    session_user_id INT PRIMARY KEY,
    last_message VARCHAR(255) NOT NULL,
    last_sent INT NOT NULL,
    FOREIGN KEY (session_user_id)
        REFERENCES chat.session_users (id)
);