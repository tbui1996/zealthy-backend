CREATE TABLE IF NOT EXISTS chat.session_notes
(
    session_id INT NOT NULL,
    notes      VARCHAR(2048),
    FOREIGN KEY (session_id)
        REFERENCES chat.sessions (id)
);