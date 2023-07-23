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
    "enumeration" TEXT NOT NULL DEFAULT '',
    "admins_only" BOOLEAN NOT NULL DEFAULT 'true',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS service_setting_configurations (
    "id" TEXT NOT NULL PRIMARY KEY,
    "value" TEXT NOT NULL DEFAULT '',
    "notes" TEXT NOT NULL DEFAULT '',
    "service_setting_id" TEXT NOT NULL REFERENCES service_settings("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "belongs_to_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "belongs_to_household" TEXT NOT NULL REFERENCES households("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("belongs_to_user", "belongs_to_household", "service_setting_id")
);
