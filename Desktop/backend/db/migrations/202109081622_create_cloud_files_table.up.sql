CREATE TABLE IF NOT EXISTS cloud.files (
  id SERIAL PRIMARY KEY,
  file_id VARCHAR(255) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  file_mimetype VARCHAR(255) NOT NULL,
  file_path VARCHAR(255) NOT NULL,

  -- TODO: Relate this to chats when Table exists
  send_user_id VARCHAR(255) NOT NULL,
  chat_id VARCHAR(255) NOT NULL,
  date_uploaded TIMESTAMP NOT NULL,
  date_last_accessed TIMESTAMP NOT NULL,
  member_id VARCHAR(255) NULL
);
