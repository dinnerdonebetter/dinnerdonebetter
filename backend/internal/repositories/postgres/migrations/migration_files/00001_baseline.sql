-- Baseline

CREATE TYPE time_zone AS ENUM (
    'UTC',
    'US/Pacific',
    'US/Mountain',
    'US/Central',
    'US/Eastern'
);

CREATE TYPE invitation_state AS ENUM (
    'pending',
    'cancelled',
    'accepted',
    'rejected'
);

CREATE TYPE setting_type AS ENUM (
    'user',
    'account',
    'membership'
);


CREATE TYPE webhook_event AS ENUM (
    'webhook_created',
    'webhook_updated',
    'webhook_archived'
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    avatar_src TEXT,
    email_address TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    password_last_changed_at TIMESTAMP WITH TIME ZONE,
    requires_password_change BOOLEAN DEFAULT FALSE NOT NULL,
    two_factor_secret TEXT NOT NULL,
    two_factor_secret_verified_at TIMESTAMP WITH TIME ZONE,
    service_role TEXT DEFAULT 'service_user'::TEXT NOT NULL,
    user_account_status TEXT DEFAULT 'unverified'::TEXT NOT NULL,
    user_account_status_explanation TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    birthday TIMESTAMP WITH TIME ZONE,
    email_address_verification_token TEXT DEFAULT ''::TEXT,
    email_address_verified_at TIMESTAMP WITH TIME ZONE,
    first_name TEXT DEFAULT ''::TEXT NOT NULL,
    last_name TEXT DEFAULT ''::TEXT NOT NULL,
    last_accepted_terms_of_service TIMESTAMP WITH TIME ZONE,
    last_accepted_privacy_policy TIMESTAMP WITH TIME ZONE,
    last_indexed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS accounts (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    billing_status TEXT DEFAULT 'unpaid'::TEXT NOT NULL,
    contact_phone TEXT DEFAULT ''::TEXT NOT NULL,
    payment_processor_customer_id TEXT DEFAULT ''::TEXT NOT NULL,
    subscription_plan_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    time_zone time_zone DEFAULT 'US/Central'::time_zone NOT NULL,
    address_line_1 TEXT DEFAULT ''::TEXT NOT NULL,
    address_line_2 TEXT DEFAULT ''::TEXT NOT NULL,
    city TEXT DEFAULT ''::TEXT NOT NULL,
    state TEXT DEFAULT ''::TEXT NOT NULL,
    zip_code TEXT DEFAULT ''::TEXT NOT NULL,
    country TEXT DEFAULT ''::TEXT NOT NULL,
    latitude NUMERIC(14,11),
    longitude NUMERIC(14,11),
    last_payment_provider_sync_occurred_at TIMESTAMP WITH TIME ZONE,
    webhook_hmac_secret TEXT DEFAULT ''::TEXT NOT NULL,
    UNIQUE(belongs_to_user, name)
);

CREATE TABLE IF NOT EXISTS account_user_memberships (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    default_account BOOLEAN DEFAULT FALSE NOT NULL,
    account_role TEXT DEFAULT 'account_user'::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_account, belongs_to_user)
);

CREATE TABLE IF NOT EXISTS account_invitations (
    id TEXT NOT NULL PRIMARY KEY,
    destination_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    to_email TEXT NOT NULL,
    to_user TEXT  REFERENCES users("id") ON DELETE CASCADE,
    from_user TEXT NOT NULL  REFERENCES users("id") ON DELETE CASCADE,
    status invitation_state DEFAULT 'pending'::invitation_state NOT NULL,
    note TEXT DEFAULT ''::TEXT NOT NULL,
    status_note TEXT DEFAULT ''::TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + '7 days'::INTERVAL) NOT NULL,
    to_name TEXT DEFAULT ''::TEXT NOT NULL,
    UNIQUE(to_user, to_email, from_user, destination_account)
);

CREATE TABLE IF NOT EXISTS service_settings (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT DEFAULT ''::TEXT NOT NULL,
    type setting_type DEFAULT 'user'::setting_type NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    default_value TEXT,
    enumeration TEXT DEFAULT ''::TEXT NOT NULL,
    admins_only BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, archived_at)
);

CREATE TABLE IF NOT EXISTS service_setting_configurations (
    id TEXT NOT NULL PRIMARY KEY,
    value TEXT DEFAULT ''::TEXT NOT NULL,
    notes TEXT DEFAULT ''::TEXT NOT NULL,
    service_setting_id TEXT NOT NULL REFERENCES service_settings("id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_user, belongs_to_account, service_setting_id)
);

CREATE TABLE IF NOT EXISTS oauth2_clients (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(client_id),
    UNIQUE(client_secret)
);

CREATE TABLE IF NOT EXISTS oauth2_client_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    client_id TEXT NOT NULL REFERENCES oauth2_clients("client_id") ON DELETE CASCADE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    redirect_uri TEXT DEFAULT ''::TEXT NOT NULL,
    code TEXT DEFAULT ''::TEXT NOT NULL,
    code_challenge TEXT DEFAULT ''::TEXT NOT NULL,
    code_challenge_method TEXT DEFAULT ''::TEXT NOT NULL,
    code_created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    code_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + '01:00:00'::INTERVAL) NOT NULL,
    access TEXT DEFAULT ''::TEXT NOT NULL,
    access_created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    access_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + '01:00:00'::INTERVAL) NOT NULL,
    refresh TEXT DEFAULT ''::TEXT NOT NULL,
    refresh_created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    refresh_expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + '01:00:00'::INTERVAL) NOT NULL,
    UNIQUE(belongs_to_user, client_id, code_expires_at, access_expires_at, refresh_expires_at)
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    redeemed_at TIMESTAMP WITH TIME ZONE,
    belongs_to_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS webhook_trigger_events (
    id TEXT NOT NULL PRIMARY KEY,
    trigger_event webhook_event NOT NULL,
    belongs_to_webhook TEXT NOT NULL REFERENCES webhooks("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(trigger_event, belongs_to_webhook)
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(token)
);

-- =============================================================================
-- INDEXES FOR BASELINE TABLES
-- =============================================================================

-- Users table indexes
CREATE INDEX idx_users_archived_at ON users (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_email_address_active ON users (email_address) WHERE archived_at IS NULL;
CREATE INDEX idx_users_username_active ON users (username) WHERE archived_at IS NULL;
CREATE INDEX idx_users_email_verification_token ON users (email_address_verification_token) WHERE archived_at IS NULL AND email_address_verification_token != '';
CREATE INDEX idx_users_service_role_username ON users (service_role, username) WHERE archived_at IS NULL;
CREATE INDEX idx_users_two_factor_verified ON users (two_factor_secret_verified_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_indexing_status ON users (last_indexed_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_active_created_at ON users (created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_active_updated_at ON users (last_updated_at) WHERE archived_at IS NULL;
CREATE INDEX idx_users_indexing_needed ON users (last_indexed_at) WHERE archived_at IS NULL;

-- Accounts table indexes
CREATE INDEX idx_accounts_belongs_to_user ON accounts (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_archived_at ON accounts (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_name ON accounts (belongs_to_user, name) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_payment_sync ON accounts (last_payment_provider_sync_occurred_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_billing_status ON accounts (billing_status) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_created_at ON accounts (belongs_to_user, created_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_updated_at ON accounts (belongs_to_user, last_updated_at) WHERE archived_at IS NULL;
CREATE INDEX idx_accounts_user_billing ON accounts (belongs_to_user, billing_status) WHERE archived_at IS NULL;

-- Account user memberships indexes
CREATE INDEX idx_memberships_user ON account_user_memberships (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_account ON account_user_memberships (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_default_account ON account_user_memberships (belongs_to_user, default_account) WHERE archived_at IS NULL AND default_account = TRUE;
CREATE INDEX idx_memberships_user_account ON account_user_memberships (belongs_to_user, belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_memberships_account_role ON account_user_memberships (belongs_to_account, account_role) WHERE archived_at IS NULL;

-- Account invitations indexes
CREATE INDEX idx_invitations_destination_account ON account_invitations (destination_account) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_from_user ON account_invitations (from_user) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_to_user ON account_invitations (to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_to_email ON account_invitations (to_email) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_token ON account_invitations (token) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_status ON account_invitations (status) WHERE archived_at IS NULL;
CREATE INDEX idx_invitations_expires_at ON account_invitations (expires_at) WHERE archived_at IS NULL;

-- Service settings indexes
CREATE INDEX idx_service_settings_archived_at ON service_settings (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_service_settings_name ON service_settings (name) WHERE archived_at IS NULL;
CREATE INDEX idx_service_settings_type ON service_settings (type) WHERE archived_at IS NULL;
CREATE INDEX idx_service_settings_admins_only ON service_settings (admins_only) WHERE archived_at IS NULL;

-- Service setting configurations indexes
CREATE INDEX idx_setting_configs_user ON service_setting_configurations (belongs_to_user) WHERE archived_at IS NULL;
CREATE INDEX idx_setting_configs_account ON service_setting_configurations (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_setting_configs_setting ON service_setting_configurations (service_setting_id) WHERE archived_at IS NULL;

-- OAuth2 clients indexes
CREATE INDEX idx_oauth2_clients_archived_at ON oauth2_clients (archived_at) WHERE archived_at IS NULL;

-- OAuth2 client tokens indexes
CREATE INDEX idx_oauth2_tokens_client_id ON oauth2_client_tokens (client_id);
CREATE INDEX idx_oauth2_tokens_user ON oauth2_client_tokens (belongs_to_user);
CREATE INDEX idx_oauth2_tokens_user_client ON oauth2_client_tokens (belongs_to_user, client_id);
CREATE INDEX idx_oauth2_tokens_code_expires ON oauth2_client_tokens (code_expires_at);
CREATE INDEX idx_oauth2_tokens_access_expires ON oauth2_client_tokens (access_expires_at);
CREATE INDEX idx_oauth2_tokens_refresh_expires ON oauth2_client_tokens (refresh_expires_at);

-- Password reset tokens indexes
CREATE INDEX idx_password_reset_user ON password_reset_tokens (belongs_to_user);
CREATE INDEX idx_password_reset_token ON password_reset_tokens (token);
CREATE INDEX idx_password_reset_expires ON password_reset_tokens (expires_at);
CREATE INDEX idx_password_reset_unredeemed ON password_reset_tokens (belongs_to_user, expires_at) WHERE redeemed_at IS NULL;

-- Webhooks indexes
CREATE INDEX idx_webhooks_account ON webhooks (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_webhooks_archived_at ON webhooks (archived_at) WHERE archived_at IS NULL;

-- Webhook trigger events indexes
CREATE INDEX idx_webhook_triggers_webhook ON webhook_trigger_events (belongs_to_webhook) WHERE archived_at IS NULL;
CREATE INDEX idx_webhook_triggers_event ON webhook_trigger_events (trigger_event) WHERE archived_at IS NULL;

-- Sessions indexes
CREATE INDEX idx_sessions_expiry ON sessions (expiry);
CREATE INDEX idx_sessions_created_at ON sessions (created_at);

-- Text search indexes (for efficient LIKE and ILIKE operations)
-- Uncomment if pg_trgm extension is available:
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX idx_users_username_trgm ON users USING gin (username gin_trgm_ops) WHERE archived_at IS NULL;
