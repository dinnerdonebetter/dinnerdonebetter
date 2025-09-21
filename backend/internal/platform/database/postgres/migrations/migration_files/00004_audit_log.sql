CREATE TYPE audit_log_event_type AS ENUM (
    'other',
    'created',
    'updated',
    'archived'
);

CREATE TABLE IF NOT EXISTS audit_log_entries (
    id TEXT NOT NULL PRIMARY KEY,
    resource_type TEXT NOT NULL,
    relevant_id TEXT NOT NULL DEFAULT '',
    event_type audit_log_event_type NOT NULL DEFAULT 'other',
    changes JSONB NOT NULL,
    belongs_to_account TEXT REFERENCES accounts("id") ON DELETE CASCADE,
    belongs_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Performance indexes for audit log (will grow large over time)
CREATE INDEX IF NOT EXISTS audit_log_account_type_created_idx ON audit_log_entries(belongs_to_account, resource_type, created_at);
CREATE INDEX IF NOT EXISTS audit_log_user_type_created_idx ON audit_log_entries(belongs_to_user, resource_type, created_at);
CREATE INDEX IF NOT EXISTS audit_log_cleanup_idx ON audit_log_entries(created_at) WHERE created_at < NOW() - INTERVAL '1 year';
