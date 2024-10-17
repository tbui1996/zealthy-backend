ALTER TABLE chat.session_statuses
ADD COLUMN opened_at TIMESTAMP;
ALTER TABLE chat.session_statuses
ADD COLUMN closed_at TIMESTAMP;
ALTER TABLE chat.session_statuses
ADD COLUMN ride_scheduled BOOLEAN;
