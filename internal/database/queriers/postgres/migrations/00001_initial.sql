CREATE TABLE sessions (
    "token" TEXT PRIMARY KEY,
    "data" BYTEA NOT NULL,
    "expiry" TIMESTAMPTZ NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
);

CREATE INDEX sessions_expiry_idx ON sessions ("expiry");

CREATE TABLE IF NOT EXISTS users (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "username" TEXT NOT NULL,
    "avatar_src" TEXT,
    "email_address" TEXT NOT NULL,
    "hashed_password" TEXT NOT NULL,
    "password_last_changed_on" INTEGER,
    "requires_password_change" BOOLEAN NOT NULL DEFAULT 'false',
    "two_factor_secret" TEXT NOT NULL,
    "two_factor_secret_verified_on" BIGINT DEFAULT NULL,
    "birth_day" SMALLINT,
    "birth_month" SMALLINT,
    "service_roles" TEXT NOT NULL DEFAULT 'service_user',
    "user_account_status" TEXT NOT NULL DEFAULT 'unverified',
    "user_account_status_explanation" TEXT NOT NULL DEFAULT '',
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("username")
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
     "id" CHAR(27) NOT NULL PRIMARY KEY,
     "token" TEXT NOT NULL,
     "expires_at" BIGINT NOT NULL,
     "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
     "last_updated_on" BIGINT DEFAULT NULL,
     "redeemed_on" BIGINT DEFAULT NULL,
     "belongs_to_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS households (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "billing_status" TEXT NOT NULL DEFAULT 'unpaid',
    "contact_email" TEXT NOT NULL DEFAULT '',
    "contact_phone" TEXT NOT NULL DEFAULT '',
    "payment_processor_customer_id" TEXT NOT NULL DEFAULT '',
    "subscription_plan_id" TEXT,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    "belongs_to_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    UNIQUE("belongs_to_user", "name")
);

CREATE INDEX IF NOT EXISTS households_belongs_to_user ON households (belongs_to_user);

CREATE TYPE invitation_state AS ENUM ('pending', 'cancelled', 'accepted', 'rejected');

CREATE TABLE IF NOT EXISTS household_invitations (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "destination_household" CHAR(27) NOT NULL REFERENCES households("id") ON DELETE CASCADE,
    "to_email" TEXT NOT NULL,
    "to_user" CHAR(27) REFERENCES users("id") ON DELETE CASCADE,
    "from_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "status" invitation_state NOT NULL DEFAULT 'pending',
    "note" TEXT NOT NULL DEFAULT '',
    "status_note" TEXT NOT NULL DEFAULT '',
    "token" TEXT NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("to_user", "from_user", "destination_household")
);

CREATE INDEX IF NOT EXISTS household_invitations_destination_household ON household_invitations (destination_household);
CREATE INDEX IF NOT EXISTS household_invitations_to_user ON household_invitations (to_user);
CREATE INDEX IF NOT EXISTS household_invitations_from_user ON household_invitations (from_user);

CREATE TABLE IF NOT EXISTS household_user_memberships (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_household" CHAR(27) NOT NULL REFERENCES households("id") ON DELETE CASCADE,
    "belongs_to_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "default_household" BOOLEAN NOT NULL DEFAULT 'false',
    "household_roles" TEXT NOT NULL DEFAULT 'household_user',
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    UNIQUE("belongs_to_household", "belongs_to_user")
);

CREATE INDEX IF NOT EXISTS household_user_memberships_belongs_to_household ON household_user_memberships (belongs_to_household);
CREATE INDEX IF NOT EXISTS household_user_memberships_belongs_to_user ON household_user_memberships (belongs_to_user);

CREATE TABLE IF NOT EXISTS api_clients (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT DEFAULT '',
    "client_id" TEXT NOT NULL,
    "secret_key" BYTEA NOT NULL,
    "permissions" BIGINT NOT NULL DEFAULT 0,
    "admin_permissions" BIGINT NOT NULL DEFAULT 0,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    "belongs_to_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS api_clients_belongs_to_user ON api_clients (belongs_to_user);

CREATE TABLE IF NOT EXISTS webhooks (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "content_type" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "method" TEXT NOT NULL,
    "events" TEXT NOT NULL,
    "data_types" TEXT NOT NULL,
    "topics" TEXT NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    "belongs_to_household" CHAR(27) NOT NULL REFERENCES households("id") ON DELETE CASCADE
);