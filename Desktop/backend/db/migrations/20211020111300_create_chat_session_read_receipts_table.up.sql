CREATE TABLE IF NOT EXISTS chat.session_read_receipts (
    session_user_id INT PRIMARY KEY,
    last_read INT NOT NULL,
    FOREIGN KEY (session_user_id)
        REFERENCES chat.session_users (id)
);