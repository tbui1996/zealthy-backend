CREATE TABLE IF NOT EXISTS chat.session_users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    session_id INT NOT NULL,
    FOREIGN KEY (session_id)
        REFERENCES chat.sessions (id)
);