CREATE TABLE IF NOT EXISTS oauth2_clients (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "client_id" TEXT NOT NULL,
    "client_secret" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("name", "client_id", "client_secret", "archived_at")
);

CREATE TYPE oauth2_client_token_scopes AS ENUM (
    'unknown',
    'household_user',
    'household_admin',
    'service_admin'
);

CREATE TABLE IF NOT EXISTS oauth2_client_tokens (
    "id" TEXT NOT NULL PRIMARY KEY,
    "client_id" TEXT NOT NULL REFERENCES oauth2_clients("id") ON DELETE CASCADE,
    "belongs_to_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "redirect_uri" TEXT NOT NULL DEFAULT '',
    "scope" oauth2_client_token_scopes NOT NULL DEFAULT 'unknown',
    "code" TEXT NOT NULL DEFAULT '',
    "code_challenge" TEXT NOT NULL DEFAULT '',
    "code_challenge_method" TEXT NOT NULL DEFAULT '',
    "code_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "code_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 week',
    "access" TEXT NOT NULL DEFAULT '',
    "access_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "access_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 week',
    "refresh" TEXT NOT NULL DEFAULT '',
    "refresh_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "refresh_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 week',
    UNIQUE("belongs_to_user", "client_id", "scope", "code_expires_at", "access_expires_at", "refresh_expires_at")
);
