CREATE TYPE user_notification_status AS ENUM (
    'unread',
    'read',
    'dismissed'
);

CREATE TABLE IF NOT EXISTS user_notifications (
    id TEXT NOT NULL PRIMARY KEY,
    content TEXT NOT NULL,
    status user_notification_status NOT NULL DEFAULT 'unread',
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE
);

-- Performance indexes for user notifications (will scale poorly without proper indexing)
CREATE INDEX IF NOT EXISTS user_notifications_user_unread_idx ON user_notifications(belongs_to_user, created_at) WHERE status = 'unread';
CREATE INDEX IF NOT EXISTS user_notifications_user_status_idx ON user_notifications(belongs_to_user, status, created_at);
