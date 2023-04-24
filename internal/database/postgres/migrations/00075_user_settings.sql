CREATE TYPE setting_type AS ENUM (
    'user',
    'household',
    'membership'
);

CREATE TABLE IF NOT EXISTS service_settings (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "type" setting_type NOT NULL DEFAULT 'user',
    "description" TEXT NOT NULL DEFAULT '',
    "default_value" TEXT,
    "admins_only" BOOLEAN NOT NULL DEFAULT 'true',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS configured_service_settings (
    "id" TEXT NOT NULL PRIMARY KEY,
    "value" TEXT NOT NULL DEFAULT '',
    "notes" TEXT NOT NULL DEFAULT '',
    "service_setting_id" TEXT NOT NULL REFERENCES service_settings("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "user_id" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "household_id" TEXT NOT NULL REFERENCES households("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("user_id", "household_id", "service_setting_id")
);
