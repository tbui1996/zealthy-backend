DROP INDEX IF EXISTS key_deleted_at_key_null, name_deleted_at_key_null;
DROP TABLE IF EXISTS feature_flags.flags;

CREATE TABLE IF NOT EXISTS feature_flags.flags (
     id SERIAL PRIMARY KEY,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
     created_by VARCHAR(100) NOT NULL,
     updated_by VARCHAR(100) NOT NULL,
     key VARCHAR(100) NOT NULL UNIQUE,
     is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
     is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
     name VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_flags_is_deleted ON feature_flags.flags (is_deleted);