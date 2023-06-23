CREATE TABLE IF NOT EXISTS oauth2_clients (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "client_id" TEXT NOT NULL,
    "client_secret" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE("client_id"),
    UNIQUE("client_secret"),
    UNIQUE("name", "archived_at")
);

INSERT INTO oauth2_clients ("id", "name", "description", "client_id", "client_secret")
VALUES (
    'ciaaq8gpbq30v6sm92p0',
    'demo',
    'Demo client',
    'a3b0030dcfc2122eec315c7f336ff0e7a89ac71565fd0cee216015d244e25bd7',
    '48cd61ca7a45e0b9a9a5cea23bf29f4b74754effe3b24ff939a00f78c78ebe47'
);

CREATE TYPE oauth2_client_token_scopes AS ENUM (
    'unknown',
    'household_member',
    'household_admin',
    'service_admin'
);

CREATE TABLE IF NOT EXISTS oauth2_client_tokens (
    "id" TEXT NOT NULL PRIMARY KEY,
    "client_id" TEXT NOT NULL REFERENCES oauth2_clients("client_id") ON DELETE CASCADE,
    "belongs_to_user" TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "redirect_uri" TEXT NOT NULL DEFAULT '',
    "scope" oauth2_client_token_scopes NOT NULL DEFAULT 'unknown',
    "code" TEXT NOT NULL DEFAULT '',
    "code_challenge" TEXT NOT NULL DEFAULT '',
    "code_challenge_method" TEXT NOT NULL DEFAULT '',
    "code_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "code_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 hour',
    "access" TEXT NOT NULL DEFAULT '',
    "access_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "access_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 hour',
    "refresh" TEXT NOT NULL DEFAULT '',
    "refresh_created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "refresh_expires_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '1 hour',
    UNIQUE("belongs_to_user", "client_id", "scope", "code_expires_at", "access_expires_at", "refresh_expires_at")
);
