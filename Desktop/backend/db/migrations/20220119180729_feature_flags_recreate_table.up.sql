DROP INDEX IF EXISTS idx_flags_is_deleted;

DROP TABLE IF EXISTS feature_flags.flags;

CREATE TABLE IF NOT EXISTS feature_flags.flags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL,
    key VARCHAR(100) NOT NULL,
    deleted_at TIMESTAMP,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT key_deleted_at_key UNIQUE (key, deleted_at),
    CONSTRAINT name_deleted_at_key UNIQUE (name, deleted_at)
);

CREATE UNIQUE INDEX IF NOT EXISTS key_deleted_at_key_null ON feature_flags.flags (key) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS name_deleted_at_key_null ON feature_flags.flags (name) WHERE deleted_at IS NULL;