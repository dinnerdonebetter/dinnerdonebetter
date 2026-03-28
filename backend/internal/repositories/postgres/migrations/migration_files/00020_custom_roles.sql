-- Custom Roles Migration
-- Admin-configurable RBAC roles with arbitrary permission subsets

CREATE TABLE IF NOT EXISTS custom_roles (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    scope TEXT NOT NULL CHECK (scope IN ('service', 'account')),
    belongs_to_account TEXT REFERENCES accounts("id") ON DELETE CASCADE,
    created_by TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (name, scope, belongs_to_account)
);

CREATE TABLE IF NOT EXISTS custom_role_permissions (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_role TEXT NOT NULL REFERENCES custom_roles("id") ON DELETE CASCADE,
    permission TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (belongs_to_role, permission)
);

CREATE TABLE IF NOT EXISTS user_custom_role_assignments (
    id TEXT NOT NULL PRIMARY KEY,
    custom_role_id TEXT NOT NULL REFERENCES custom_roles("id") ON DELETE CASCADE,
    user_id TEXT REFERENCES users("id") ON DELETE CASCADE,
    account_membership_id TEXT REFERENCES account_user_memberships("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (custom_role_id, user_id, account_membership_id),
    CHECK (
        (user_id IS NOT NULL AND account_membership_id IS NULL)
        OR (user_id IS NULL AND account_membership_id IS NOT NULL)
    )
);

CREATE INDEX idx_custom_roles_scope ON custom_roles (scope) WHERE archived_at IS NULL;
CREATE INDEX idx_custom_roles_account ON custom_roles (belongs_to_account) WHERE archived_at IS NULL;
CREATE INDEX idx_custom_role_permissions_role ON custom_role_permissions (belongs_to_role);
CREATE INDEX idx_user_custom_role_assignments_user ON user_custom_role_assignments (user_id) WHERE archived_at IS NULL;
CREATE INDEX idx_user_custom_role_assignments_membership ON user_custom_role_assignments (account_membership_id) WHERE archived_at IS NULL;
