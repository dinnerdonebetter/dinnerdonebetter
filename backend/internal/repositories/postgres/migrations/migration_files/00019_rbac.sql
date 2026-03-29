-- RBAC Migration
-- Database-driven roles, permissions, role-permission mappings, role hierarchy, and user-role assignments

-- Drop old hardcoded role columns (no deployed DB, no data to migrate)
ALTER TABLE users DROP COLUMN IF EXISTS service_role;
ALTER TABLE account_user_memberships DROP COLUMN IF EXISTS account_role;

-- =============================================================================
-- TABLE: user_roles
-- =============================================================================

CREATE TABLE IF NOT EXISTS user_roles (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    scope TEXT NOT NULL CHECK (scope IN ('service', 'account')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name)
);

CREATE INDEX idx_user_roles_archived ON user_roles (archived_at) WHERE archived_at IS NULL;

-- =============================================================================
-- TABLE: permissions
-- =============================================================================

CREATE TABLE IF NOT EXISTS permissions (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''::TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name)
);

CREATE INDEX idx_permissions_archived ON permissions (archived_at) WHERE archived_at IS NULL;
CREATE INDEX idx_permissions_name ON permissions (name) WHERE archived_at IS NULL;

-- =============================================================================
-- TABLE: user_role_permissions (many-to-many role<->permission)
-- =============================================================================

CREATE TABLE IF NOT EXISTS user_role_permissions (
    id TEXT NOT NULL PRIMARY KEY,
    role_id TEXT NOT NULL REFERENCES user_roles("id") ON DELETE CASCADE,
    permission_id TEXT NOT NULL REFERENCES permissions("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_user_role_permissions_role ON user_role_permissions (role_id) WHERE archived_at IS NULL;
CREATE INDEX idx_user_role_permissions_permission ON user_role_permissions (permission_id) WHERE archived_at IS NULL;

-- =============================================================================
-- TABLE: user_role_assignments (maps users to roles, optional account scope)
-- =============================================================================

CREATE TABLE IF NOT EXISTS user_role_assignments (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    role_id TEXT NOT NULL REFERENCES user_roles("id") ON DELETE CASCADE,
    account_id TEXT REFERENCES accounts("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_user_role_assignments_unique_with_account
    ON user_role_assignments (user_id, role_id, account_id)
    WHERE archived_at IS NULL AND account_id IS NOT NULL;

CREATE UNIQUE INDEX idx_user_role_assignments_unique_without_account
    ON user_role_assignments (user_id, role_id)
    WHERE archived_at IS NULL AND account_id IS NULL;

CREATE INDEX idx_user_role_assignments_user ON user_role_assignments (user_id) WHERE archived_at IS NULL;
CREATE INDEX idx_user_role_assignments_user_account ON user_role_assignments (user_id, account_id) WHERE archived_at IS NULL;
CREATE INDEX idx_user_role_assignments_account ON user_role_assignments (account_id) WHERE archived_at IS NULL AND account_id IS NOT NULL;

-- =============================================================================
-- TABLE: user_role_hierarchy (child inherits parent's permissions)
-- =============================================================================

CREATE TABLE IF NOT EXISTS user_role_hierarchy (
    id TEXT NOT NULL PRIMARY KEY,
    parent_role_id TEXT NOT NULL REFERENCES user_roles("id") ON DELETE CASCADE,
    child_role_id TEXT NOT NULL REFERENCES user_roles("id") ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(parent_role_id, child_role_id),
    CHECK(parent_role_id != child_role_id)
);

CREATE INDEX idx_user_role_hierarchy_child ON user_role_hierarchy (child_role_id) WHERE archived_at IS NULL;
CREATE INDEX idx_user_role_hierarchy_parent ON user_role_hierarchy (parent_role_id) WHERE archived_at IS NULL;

-- =============================================================================
-- SEED DATA: roles
-- =============================================================================

INSERT INTO user_roles (id, name, description, scope) VALUES
    ('role_service_user', 'service_user', 'Default non-admin user', 'service'),
    ('role_service_admin', 'service_admin', 'Full service administrator', 'service'),
    ('role_service_data_admin', 'service_data_admin', 'Service data management', 'service'),
    ('role_account_member', 'account_member', 'Basic account member', 'account'),
    ('role_account_admin', 'account_admin', 'Account administrator', 'account');

-- =============================================================================
-- SEED DATA: role hierarchy
-- =============================================================================

INSERT INTO user_role_hierarchy (id, parent_role_id, child_role_id) VALUES
    ('d74elbkn9qd4nv86eo20', 'role_account_member', 'role_account_admin'),
    ('d74elbkn9qd4nv86eo2g', 'role_account_admin', 'role_service_admin');

-- =============================================================================
-- SEED DATA: core permissions (non-mealplanning)
-- =============================================================================

INSERT INTO permissions (id, name, description) VALUES
    -- admin
    ('perm_admin_read_user_data', 'admin.read_user_data', 'Read user data (admin)'),
    ('perm_queues_publish_message', 'queues.publish.message', 'Publish arbitrary queue messages'),
    -- analytics
    ('perm_report_analytics_events', 'report.analytics_events', 'Report analytics events'),
    -- auth
    ('perm_update_user_status', 'update.user_status', 'Update user account status'),
    ('perm_imitate_user', 'imitate.user', 'Impersonate users'),
    ('perm_read_user', 'read.user', 'Read user data'),
    ('perm_search_user', 'search.user', 'Search users'),
    ('perm_archive_user', 'archive.user', 'Archive users'),
    ('perm_manage_user_sessions', 'manage.user_sessions', 'Manage user sessions'),
    -- identity
    ('perm_update_account', 'update.account', 'Update account'),
    ('perm_archive_account', 'archive.account', 'Archive account'),
    ('perm_transfer_account', 'transfer.account', 'Transfer account ownership'),
    ('perm_account_add_member', 'account.add.member', 'Invite user to account'),
    ('perm_account_membership_modify', 'account.membership.modify', 'Modify member permissions'),
    ('perm_remove_member_account', 'remove_member.account', 'Remove member from account'),
    -- oauth
    ('perm_create_oauth2_clients', 'create.oauth2_clients', 'Create OAuth2 clients'),
    ('perm_read_oauth2_clients', 'read.oauth2_clients', 'Read OAuth2 clients'),
    ('perm_archive_oauth2_clients', 'archive.oauth2_clients', 'Archive OAuth2 clients'),
    -- settings
    ('perm_create_service_settings', 'create.service_settings', 'Create service settings'),
    ('perm_read_service_settings', 'read.service_settings', 'Read service settings'),
    ('perm_search_service_settings', 'search.service_settings', 'Search service settings'),
    ('perm_archive_service_settings', 'archive.service_settings', 'Archive service settings'),
    ('perm_create_service_setting_configurations', 'create.service_setting_configurations', 'Create service setting configurations'),
    ('perm_read_service_setting_configurations', 'read.service_setting_configurations', 'Read service setting configurations'),
    ('perm_update_service_setting_configurations', 'update.service_setting_configurations', 'Update service setting configurations'),
    ('perm_archive_service_setting_configurations', 'archive.service_setting_configurations', 'Archive service setting configurations'),
    -- notifications
    ('perm_create_user_notifications', 'create.user_notifications', 'Create user notifications'),
    ('perm_read_user_notifications', 'read.user_notifications', 'Read user notifications'),
    ('perm_update_user_notifications', 'update.user_notifications', 'Update user notifications'),
    ('perm_create_user_device_tokens', 'create.user_device_tokens', 'Create user device tokens'),
    ('perm_read_user_device_tokens', 'read.user_device_tokens', 'Read user device tokens'),
    ('perm_archive_user_device_tokens', 'archive.user_device_tokens', 'Archive user device tokens'),
    -- webhooks
    ('perm_create_webhooks', 'create.webhooks', 'Create webhooks'),
    ('perm_read_webhooks', 'read.webhooks', 'Read webhooks'),
    ('perm_update_webhooks', 'update.webhooks', 'Update webhooks'),
    ('perm_archive_webhooks', 'archive.webhooks', 'Archive webhooks'),
    ('perm_create_webhook_trigger_configs', 'create.webhook_trigger_configs', 'Create webhook trigger configs'),
    ('perm_archive_webhook_trigger_configs', 'archive.webhook_trigger_configs', 'Archive webhook trigger configs'),
    ('perm_create_webhook_trigger_events', 'create.webhook_trigger_events', 'Create webhook trigger events'),
    ('perm_read_webhook_trigger_events', 'read.webhook_trigger_events', 'Read webhook trigger events'),
    ('perm_update_webhook_trigger_events', 'update.webhook_trigger_events', 'Update webhook trigger events'),
    ('perm_archive_webhook_trigger_events', 'archive.webhook_trigger_events', 'Archive webhook trigger events'),
    -- uploaded media
    ('perm_create_uploaded_media', 'create.uploaded_media', 'Create uploaded media'),
    ('perm_read_uploaded_media', 'read.uploaded_media', 'Read uploaded media'),
    ('perm_update_uploaded_media', 'update.uploaded_media', 'Update uploaded media'),
    ('perm_archive_uploaded_media', 'archive.uploaded_media', 'Archive uploaded media'),
    -- issue reports
    ('perm_create_issue_reports', 'create.issue_reports', 'Create issue reports'),
    ('perm_read_issue_reports', 'read.issue_reports', 'Read issue reports'),
    ('perm_update_issue_reports', 'update.issue_reports', 'Update issue reports'),
    ('perm_archive_issue_reports', 'archive.issue_reports', 'Archive issue reports'),
    -- audit
    ('perm_read_audit_log_entries', 'read.audit_log_entries', 'Read audit log entries'),
    -- comments
    ('perm_create_comments', 'create.comments', 'Create comments'),
    ('perm_read_comments', 'read.comments', 'Read comments'),
    ('perm_update_comments', 'update.comments', 'Update comments'),
    ('perm_archive_comments', 'archive.comments', 'Archive comments'),
    -- waitlists
    ('perm_create_waitlists', 'create.waitlists', 'Create waitlists'),
    ('perm_read_waitlists', 'read.waitlists', 'Read waitlists'),
    ('perm_update_waitlists', 'update.waitlists', 'Update waitlists'),
    ('perm_archive_waitlists', 'archive.waitlists', 'Archive waitlists'),
    ('perm_create_waitlist_signups', 'create.waitlist_signups', 'Create waitlist signups'),
    ('perm_read_waitlist_signups', 'read.waitlist_signups', 'Read waitlist signups'),
    ('perm_update_waitlist_signups', 'update.waitlist_signups', 'Update waitlist signups'),
    ('perm_archive_waitlist_signups', 'archive.waitlist_signups', 'Archive waitlist signups'),
    -- payments
    ('perm_create_checkout_sessions', 'create.checkout_sessions', 'Create checkout sessions'),
    ('perm_cancel_subscriptions', 'cancel.subscriptions', 'Cancel subscriptions'),
    ('perm_read_purchases', 'read.purchases', 'Read purchases'),
    ('perm_read_payment_history', 'read.payment_history', 'Read payment history'),
    ('perm_create_products', 'create.products', 'Create products'),
    ('perm_read_products', 'read.products', 'Read products'),
    ('perm_update_products', 'update.products', 'Update products'),
    ('perm_archive_products', 'archive.products', 'Archive products'),
    ('perm_create_subscriptions', 'create.subscriptions', 'Create subscriptions'),
    ('perm_read_subscriptions', 'read.subscriptions', 'Read subscriptions'),
    ('perm_update_subscriptions', 'update.subscriptions', 'Update subscriptions'),
    ('perm_archive_subscriptions', 'archive.subscriptions', 'Archive subscriptions');

-- =============================================================================
-- SEED DATA: core role-permission mappings
-- =============================================================================

-- service_admin: core permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_sa_1', 'role_service_admin', 'perm_admin_read_user_data'),
    ('urp_sa_2', 'role_service_admin', 'perm_update_user_status'),
    ('urp_sa_3', 'role_service_admin', 'perm_read_user'),
    ('urp_sa_4', 'role_service_admin', 'perm_search_user'),
    ('urp_sa_5', 'role_service_admin', 'perm_archive_user'),
    ('urp_sa_6', 'role_service_admin', 'perm_create_oauth2_clients'),
    ('urp_sa_7', 'role_service_admin', 'perm_archive_oauth2_clients'),
    ('urp_sa_8', 'role_service_admin', 'perm_archive_service_settings'),
    ('urp_sa_9', 'role_service_admin', 'perm_create_user_notifications'),
    ('urp_sa_10', 'role_service_admin', 'perm_imitate_user'),
    ('urp_sa_11', 'role_service_admin', 'perm_manage_user_sessions'),
    ('urp_sa_12', 'role_service_admin', 'perm_queues_publish_message'),
    ('urp_sa_14', 'role_service_admin', 'perm_create_service_settings'),
    ('urp_sa_17', 'role_service_admin', 'perm_create_waitlists'),
    ('urp_sa_18', 'role_service_admin', 'perm_update_waitlists'),
    ('urp_sa_19', 'role_service_admin', 'perm_archive_waitlists'),
    ('urp_sa_20', 'role_service_admin', 'perm_create_products'),
    ('urp_sa_21', 'role_service_admin', 'perm_read_products'),
    ('urp_sa_22', 'role_service_admin', 'perm_update_products'),
    ('urp_sa_23', 'role_service_admin', 'perm_archive_products'),
    ('urp_sa_24', 'role_service_admin', 'perm_create_subscriptions'),
    ('urp_sa_25', 'role_service_admin', 'perm_read_subscriptions'),
    ('urp_sa_26', 'role_service_admin', 'perm_update_subscriptions'),
    ('urp_sa_27', 'role_service_admin', 'perm_archive_subscriptions');

-- account_admin: core permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_aa_1', 'role_account_admin', 'perm_update_account'),
    ('urp_aa_2', 'role_account_admin', 'perm_archive_account'),
    ('urp_aa_3', 'role_account_admin', 'perm_transfer_account'),
    ('urp_aa_4', 'role_account_admin', 'perm_account_add_member'),
    ('urp_aa_5', 'role_account_admin', 'perm_account_membership_modify'),
    ('urp_aa_6', 'role_account_admin', 'perm_remove_member_account'),
    ('urp_aa_7', 'role_account_admin', 'perm_create_webhooks'),
    ('urp_aa_8', 'role_account_admin', 'perm_update_webhooks'),
    ('urp_aa_9', 'role_account_admin', 'perm_archive_webhooks'),
    ('urp_aa_10', 'role_account_admin', 'perm_create_issue_reports'),
    ('urp_aa_11', 'role_account_admin', 'perm_update_issue_reports'),
    ('urp_aa_12', 'role_account_admin', 'perm_archive_issue_reports'),
    ('urp_aa_25', 'role_account_admin', 'perm_create_webhook_trigger_configs'),
    ('urp_aa_26', 'role_account_admin', 'perm_archive_webhook_trigger_configs'),
    ('urp_aa_27', 'role_account_admin', 'perm_create_webhook_trigger_events'),
    ('urp_aa_28', 'role_account_admin', 'perm_read_webhook_trigger_events'),
    ('urp_aa_29', 'role_account_admin', 'perm_update_webhook_trigger_events'),
    ('urp_aa_30', 'role_account_admin', 'perm_archive_webhook_trigger_events'),
    ('urp_aa_39', 'role_account_admin', 'perm_create_checkout_sessions'),
    ('urp_aa_40', 'role_account_admin', 'perm_cancel_subscriptions'),
    ('urp_aa_41', 'role_account_admin', 'perm_read_purchases'),
    ('urp_aa_42', 'role_account_admin', 'perm_read_payment_history'),
    ('urp_aa_43', 'role_account_admin', 'perm_read_subscriptions');

-- account_member: core permissions
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_am_1', 'role_account_member', 'perm_report_analytics_events'),
    ('urp_am_2', 'role_account_member', 'perm_read_webhooks'),
    ('urp_am_3', 'role_account_member', 'perm_read_issue_reports'),
    ('urp_am_4', 'role_account_member', 'perm_read_audit_log_entries'),
    ('urp_am_5', 'role_account_member', 'perm_read_oauth2_clients'),
    ('urp_am_6', 'role_account_member', 'perm_read_service_settings'),
    ('urp_am_7', 'role_account_member', 'perm_search_service_settings'),
    ('urp_am_8', 'role_account_member', 'perm_create_uploaded_media'),
    ('urp_am_9', 'role_account_member', 'perm_read_uploaded_media'),
    ('urp_am_10', 'role_account_member', 'perm_update_uploaded_media'),
    ('urp_am_11', 'role_account_member', 'perm_archive_uploaded_media'),
    ('urp_am_96', 'role_account_member', 'perm_create_service_setting_configurations'),
    ('urp_am_97', 'role_account_member', 'perm_read_service_setting_configurations'),
    ('urp_am_98', 'role_account_member', 'perm_update_service_setting_configurations'),
    ('urp_am_99', 'role_account_member', 'perm_archive_service_setting_configurations'),
    ('urp_am_109', 'role_account_member', 'perm_create_comments'),
    ('urp_am_110', 'role_account_member', 'perm_read_comments'),
    ('urp_am_111', 'role_account_member', 'perm_update_comments'),
    ('urp_am_112', 'role_account_member', 'perm_archive_comments'),
    ('urp_am_115', 'role_account_member', 'perm_read_user_notifications'),
    ('urp_am_116', 'role_account_member', 'perm_update_user_notifications'),
    ('urp_am_117', 'role_account_member', 'perm_create_user_device_tokens'),
    ('urp_am_118', 'role_account_member', 'perm_read_user_device_tokens'),
    ('urp_am_119', 'role_account_member', 'perm_archive_user_device_tokens'),
    ('urp_am_120', 'role_account_member', 'perm_create_waitlist_signups'),
    ('urp_am_121', 'role_account_member', 'perm_update_waitlist_signups'),
    ('urp_am_122', 'role_account_member', 'perm_archive_waitlist_signups'),
    ('urp_am_123', 'role_account_member', 'perm_read_waitlists'),
    ('urp_am_124', 'role_account_member', 'perm_read_waitlist_signups'),
    ('urp_am_126', 'role_account_member', 'perm_create_checkout_sessions'),
    ('urp_am_127', 'role_account_member', 'perm_cancel_subscriptions'),
    ('urp_am_128', 'role_account_member', 'perm_read_purchases'),
    ('urp_am_129', 'role_account_member', 'perm_read_payment_history'),
    ('urp_am_130', 'role_account_member', 'perm_read_subscriptions');
