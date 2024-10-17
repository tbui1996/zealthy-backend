ALTER TABLE chat.session_statuses
DROP COLUMN IF EXISTS opened_at;
ALTER TABLE chat.session_statuses
DROP COLUMN IF EXISTS closed_at;
ALTER TABLE chat.session_statuses
DROP COLUMN IF EXISTS ride_scheduled;
