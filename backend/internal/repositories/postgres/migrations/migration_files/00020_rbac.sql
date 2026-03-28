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
-- account_admin inherits account_member
-- service_admin inherits account_admin (and transitively account_member)
-- =============================================================================

INSERT INTO user_role_hierarchy (id, parent_role_id, child_role_id) VALUES
    ('rh_acctadmin_acctmember', 'role_account_member', 'role_account_admin'),
    ('rh_svcadmin_acctadmin', 'role_account_admin', 'role_service_admin');

-- =============================================================================
-- SEED DATA: permissions
-- =============================================================================

INSERT INTO permissions (id, name, description) VALUES
    -- admin permissions
    ('perm_admin_read_user_data', 'admin.read_user_data', 'Read user data (admin)'),
    ('perm_queues_publish_message', 'queues.publish.message', 'Publish arbitrary queue messages'),
    -- analytics permissions
    ('perm_report_analytics_events', 'report.analytics_events', 'Report analytics events'),
    -- auth permissions
    ('perm_update_user_status', 'update.user_status', 'Update user account status'),
    ('perm_imitate_user', 'imitate.user', 'Impersonate users'),
    ('perm_read_user', 'read.user', 'Read user data'),
    ('perm_search_user', 'search.user', 'Search users'),
    ('perm_archive_user', 'archive.user', 'Archive users'),
    ('perm_manage_user_sessions', 'manage.user_sessions', 'Manage user sessions'),
    -- identity permissions
    ('perm_update_account', 'update.account', 'Update account'),
    ('perm_archive_account', 'archive.account', 'Archive account'),
    ('perm_transfer_account', 'transfer.account', 'Transfer account ownership'),
    ('perm_account_add_member', 'account.add.member', 'Invite user to account'),
    ('perm_account_membership_modify', 'account.membership.modify', 'Modify member permissions'),
    ('perm_remove_member_account', 'remove_member.account', 'Remove member from account'),
    -- oauth permissions
    ('perm_create_oauth2_clients', 'create.oauth2_clients', 'Create OAuth2 clients'),
    ('perm_read_oauth2_clients', 'read.oauth2_clients', 'Read OAuth2 clients'),
    ('perm_archive_oauth2_clients', 'archive.oauth2_clients', 'Archive OAuth2 clients'),
    -- settings permissions
    ('perm_create_service_settings', 'create.service_settings', 'Create service settings'),
    ('perm_read_service_settings', 'read.service_settings', 'Read service settings'),
    ('perm_search_service_settings', 'search.service_settings', 'Search service settings'),
    ('perm_archive_service_settings', 'archive.service_settings', 'Archive service settings'),
    ('perm_create_service_setting_configurations', 'create.service_setting_configurations', 'Create service setting configurations'),
    ('perm_read_service_setting_configurations', 'read.service_setting_configurations', 'Read service setting configurations'),
    ('perm_update_service_setting_configurations', 'update.service_setting_configurations', 'Update service setting configurations'),
    ('perm_archive_service_setting_configurations', 'archive.service_setting_configurations', 'Archive service setting configurations'),
    -- notifications permissions
    ('perm_create_user_notifications', 'create.user_notifications', 'Create user notifications'),
    ('perm_read_user_notifications', 'read.user_notifications', 'Read user notifications'),
    ('perm_update_user_notifications', 'update.user_notifications', 'Update user notifications'),
    ('perm_create_user_device_tokens', 'create.user_device_tokens', 'Create user device tokens'),
    ('perm_read_user_device_tokens', 'read.user_device_tokens', 'Read user device tokens'),
    ('perm_archive_user_device_tokens', 'archive.user_device_tokens', 'Archive user device tokens'),
    -- webhooks permissions
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
    -- uploaded media permissions
    ('perm_create_uploaded_media', 'create.uploaded_media', 'Create uploaded media'),
    ('perm_read_uploaded_media', 'read.uploaded_media', 'Read uploaded media'),
    ('perm_update_uploaded_media', 'update.uploaded_media', 'Update uploaded media'),
    ('perm_archive_uploaded_media', 'archive.uploaded_media', 'Archive uploaded media'),
    -- issue reports permissions
    ('perm_create_issue_reports', 'create.issue_reports', 'Create issue reports'),
    ('perm_read_issue_reports', 'read.issue_reports', 'Read issue reports'),
    ('perm_update_issue_reports', 'update.issue_reports', 'Update issue reports'),
    ('perm_archive_issue_reports', 'archive.issue_reports', 'Archive issue reports'),
    -- audit permissions
    ('perm_read_audit_log_entries', 'read.audit_log_entries', 'Read audit log entries'),
    -- comments permissions
    ('perm_create_comments', 'create.comments', 'Create comments'),
    ('perm_read_comments', 'read.comments', 'Read comments'),
    ('perm_update_comments', 'update.comments', 'Update comments'),
    ('perm_archive_comments', 'archive.comments', 'Archive comments'),
    -- waitlists permissions
    ('perm_create_waitlists', 'create.waitlists', 'Create waitlists'),
    ('perm_read_waitlists', 'read.waitlists', 'Read waitlists'),
    ('perm_update_waitlists', 'update.waitlists', 'Update waitlists'),
    ('perm_archive_waitlists', 'archive.waitlists', 'Archive waitlists'),
    ('perm_create_waitlist_signups', 'create.waitlist_signups', 'Create waitlist signups'),
    ('perm_read_waitlist_signups', 'read.waitlist_signups', 'Read waitlist signups'),
    ('perm_update_waitlist_signups', 'update.waitlist_signups', 'Update waitlist signups'),
    ('perm_archive_waitlist_signups', 'archive.waitlist_signups', 'Archive waitlist signups'),
    -- payments permissions
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
    ('perm_archive_subscriptions', 'archive.subscriptions', 'Archive subscriptions'),
    -- meal planning: valid instruments
    ('perm_create_valid_instruments', 'create.valid_instruments', 'Create valid instruments'),
    ('perm_read_valid_instruments', 'read.valid_instruments', 'Read valid instruments'),
    ('perm_search_valid_instruments', 'search.valid_instruments', 'Search valid instruments'),
    ('perm_update_valid_instruments', 'update.valid_instruments', 'Update valid instruments'),
    ('perm_archive_valid_instruments', 'archive.valid_instruments', 'Archive valid instruments'),
    -- meal planning: valid vessels
    ('perm_create_valid_vessels', 'create.valid_vessels', 'Create valid vessels'),
    ('perm_read_valid_vessels', 'read.valid_vessels', 'Read valid vessels'),
    ('perm_search_valid_vessels', 'search.valid_vessels', 'Search valid vessels'),
    ('perm_update_valid_vessels', 'update.valid_vessels', 'Update valid vessels'),
    ('perm_archive_valid_vessels', 'archive.valid_vessels', 'Archive valid vessels'),
    -- meal planning: valid ingredients
    ('perm_create_valid_ingredients', 'create.valid_ingredients', 'Create valid ingredients'),
    ('perm_read_valid_ingredients', 'read.valid_ingredients', 'Read valid ingredients'),
    ('perm_search_valid_ingredients', 'search.valid_ingredients', 'Search valid ingredients'),
    ('perm_update_valid_ingredients', 'update.valid_ingredients', 'Update valid ingredients'),
    ('perm_archive_valid_ingredients', 'archive.valid_ingredients', 'Archive valid ingredients'),
    -- meal planning: valid ingredient groups
    ('perm_create_valid_ingredient_groups', 'create.valid_ingredient_groups', 'Create valid ingredient groups'),
    ('perm_read_valid_ingredient_groups', 'read.valid_ingredient_groups', 'Read valid ingredient groups'),
    ('perm_search_valid_ingredient_groups', 'search.valid_ingredient_groups', 'Search valid ingredient groups'),
    ('perm_update_valid_ingredient_groups', 'update.valid_ingredient_groups', 'Update valid ingredient groups'),
    ('perm_archive_valid_ingredient_groups', 'archive.valid_ingredient_groups', 'Archive valid ingredient groups'),
    -- meal planning: valid preparations
    ('perm_create_valid_preparations', 'create.valid_preparations', 'Create valid preparations'),
    ('perm_read_valid_preparations', 'read.valid_preparations', 'Read valid preparations'),
    ('perm_search_valid_preparations', 'search.valid_preparations', 'Search valid preparations'),
    ('perm_update_valid_preparations', 'update.valid_preparations', 'Update valid preparations'),
    ('perm_archive_valid_preparations', 'archive.valid_preparations', 'Archive valid preparations'),
    -- meal planning: measurement units
    ('perm_create_measurement_units', 'create.measurement_units', 'Create measurement units'),
    ('perm_read_measurement_units', 'read.measurement_units', 'Read measurement units'),
    ('perm_search_measurement_units', 'search.measurement_units', 'Search measurement units'),
    ('perm_update_measurement_units', 'update.measurement_units', 'Update measurement units'),
    ('perm_archive_measurement_units', 'archive.measurement_units', 'Archive measurement units'),
    -- meal planning: measurement conversions
    ('perm_create_measurement_conversions', 'create.measurement_conversions', 'Create measurement conversions'),
    ('perm_read_measurement_conversions', 'read.measurement_conversions', 'Read measurement conversions'),
    ('perm_update_measurement_conversions', 'update.measurement_conversions', 'Update measurement conversions'),
    ('perm_archive_measurement_conversions', 'archive.measurement_conversions', 'Archive measurement conversions'),
    -- meal planning: valid ingredient preparations
    ('perm_create_valid_ingredient_preparations', 'create.valid_ingredient_preparations', 'Create valid ingredient preparations'),
    ('perm_read_valid_ingredient_preparations', 'read.valid_ingredient_preparations', 'Read valid ingredient preparations'),
    ('perm_search_valid_ingredient_preparations', 'search.valid_ingredient_preparations', 'Search valid ingredient preparations'),
    ('perm_update_valid_ingredient_preparations', 'update.valid_ingredient_preparations', 'Update valid ingredient preparations'),
    ('perm_archive_valid_ingredient_preparations', 'archive.valid_ingredient_preparations', 'Archive valid ingredient preparations'),
    -- meal planning: valid prep task configs
    ('perm_create_valid_prep_task_configs', 'create.valid_prep_task_configs', 'Create valid prep task configs'),
    ('perm_read_valid_prep_task_configs', 'read.valid_prep_task_configs', 'Read valid prep task configs'),
    ('perm_update_valid_prep_task_configs', 'update.valid_prep_task_configs', 'Update valid prep task configs'),
    ('perm_archive_valid_prep_task_configs', 'archive.valid_prep_task_configs', 'Archive valid prep task configs'),
    -- meal planning: valid ingredient state ingredients
    ('perm_create_valid_ingredient_state_ingredients', 'create.valid_ingredient_state_ingredients', 'Create valid ingredient state ingredients'),
    ('perm_read_valid_ingredient_state_ingredients', 'read.valid_ingredient_state_ingredients', 'Read valid ingredient state ingredients'),
    ('perm_search_valid_ingredient_state_ingredients', 'search.valid_ingredient_state_ingredients', 'Search valid ingredient state ingredients'),
    ('perm_update_valid_ingredient_state_ingredients', 'update.valid_ingredient_state_ingredients', 'Update valid ingredient state ingredients'),
    ('perm_archive_valid_ingredient_state_ingredients', 'archive.valid_ingredient_state_ingredients', 'Archive valid ingredient state ingredients'),
    -- meal planning: valid preparation instruments
    ('perm_create_valid_preparation_instruments', 'create.valid_preparation_instruments', 'Create valid preparation instruments'),
    ('perm_read_valid_preparation_instruments', 'read.valid_preparation_instruments', 'Read valid preparation instruments'),
    ('perm_search_valid_preparation_instruments', 'search.valid_preparation_instruments', 'Search valid preparation instruments'),
    ('perm_update_valid_preparation_instruments', 'update.valid_preparation_instruments', 'Update valid preparation instruments'),
    ('perm_archive_valid_preparation_instruments', 'archive.valid_preparation_instruments', 'Archive valid preparation instruments'),
    -- meal planning: valid preparation vessels
    ('perm_create_valid_preparation_vessels', 'create.valid_preparation_vessels', 'Create valid preparation vessels'),
    ('perm_read_valid_preparation_vessels', 'read.valid_preparation_vessels', 'Read valid preparation vessels'),
    ('perm_search_valid_preparation_vessels', 'search.valid_preparation_vessels', 'Search valid preparation vessels'),
    ('perm_update_valid_preparation_vessels', 'update.valid_preparation_vessels', 'Update valid preparation vessels'),
    ('perm_archive_valid_preparation_vessels', 'archive.valid_preparation_vessels', 'Archive valid preparation vessels'),
    -- meal planning: valid ingredient measurement units
    ('perm_create_valid_ingredient_measurement_units', 'create.valid_ingredient_measurement_units', 'Create valid ingredient measurement units'),
    ('perm_read_valid_ingredient_measurement_units', 'read.valid_ingredient_measurement_units', 'Read valid ingredient measurement units'),
    ('perm_search_valid_ingredient_measurement_units', 'search.valid_ingredient_measurement_units', 'Search valid ingredient measurement units'),
    ('perm_update_valid_ingredient_measurement_units', 'update.valid_ingredient_measurement_units', 'Update valid ingredient measurement units'),
    ('perm_archive_valid_ingredient_measurement_units', 'archive.valid_ingredient_measurement_units', 'Archive valid ingredient measurement units'),
    -- meal planning: valid ingredient states
    ('perm_create_valid_ingredient_states', 'create.valid_ingredient_states', 'Create valid ingredient states'),
    ('perm_read_valid_ingredient_states', 'read.valid_ingredient_states', 'Read valid ingredient states'),
    ('perm_update_valid_ingredient_states', 'update.valid_ingredient_states', 'Update valid ingredient states'),
    ('perm_archive_valid_ingredient_states', 'archive.valid_ingredient_states', 'Archive valid ingredient states'),
    -- meal planning: meals
    ('perm_create_meals', 'create.meals', 'Create meals'),
    ('perm_read_meals', 'read.meals', 'Read meals'),
    ('perm_update_meals', 'update.meals', 'Update meals'),
    ('perm_archive_meals', 'archive.meals', 'Archive meals'),
    -- meal planning: recipes
    ('perm_create_recipes', 'create.recipes', 'Create recipes'),
    ('perm_read_recipes', 'read.recipes', 'Read recipes'),
    ('perm_search_recipes', 'search.recipes', 'Search recipes'),
    ('perm_update_recipes', 'update.recipes', 'Update recipes'),
    ('perm_archive_recipes', 'archive.recipes', 'Archive recipes'),
    ('perm_update_recipe_status', 'update.recipe_status', 'Update recipe status'),
    -- meal planning: recipe steps
    ('perm_create_recipe_steps', 'create.recipe_steps', 'Create recipe steps'),
    ('perm_read_recipe_steps', 'read.recipe_steps', 'Read recipe steps'),
    ('perm_search_recipe_steps', 'search.recipe_steps', 'Search recipe steps'),
    ('perm_update_recipe_steps', 'update.recipe_steps', 'Update recipe steps'),
    ('perm_archive_recipe_steps', 'archive.recipe_steps', 'Archive recipe steps'),
    -- meal planning: recipe prep tasks
    ('perm_create_recipe_prep_tasks', 'create.recipe_prep_tasks', 'Create recipe prep tasks'),
    ('perm_read_recipe_prep_tasks', 'read.recipe_prep_tasks', 'Read recipe prep tasks'),
    ('perm_update_recipe_prep_tasks', 'update.recipe_prep_tasks', 'Update recipe prep tasks'),
    ('perm_archive_recipe_prep_tasks', 'archive.recipe_prep_tasks', 'Archive recipe prep tasks'),
    -- meal planning: recipe step instruments
    ('perm_create_recipe_step_instruments', 'create.recipe_step_instruments', 'Create recipe step instruments'),
    ('perm_read_recipe_step_instruments', 'read.recipe_step_instruments', 'Read recipe step instruments'),
    ('perm_search_recipe_step_instruments', 'search.recipe_step_instruments', 'Search recipe step instruments'),
    ('perm_update_recipe_step_instruments', 'update.recipe_step_instruments', 'Update recipe step instruments'),
    ('perm_archive_recipe_step_instruments', 'archive.recipe_step_instruments', 'Archive recipe step instruments'),
    -- meal planning: recipe step vessels
    ('perm_create_recipe_step_vessels', 'create.recipe_step_vessels', 'Create recipe step vessels'),
    ('perm_read_recipe_step_vessels', 'read.recipe_step_vessels', 'Read recipe step vessels'),
    ('perm_search_recipe_step_vessels', 'search.recipe_step_vessels', 'Search recipe step vessels'),
    ('perm_update_recipe_step_vessels', 'update.recipe_step_vessels', 'Update recipe step vessels'),
    ('perm_archive_recipe_step_vessels', 'archive.recipe_step_vessels', 'Archive recipe step vessels'),
    -- meal planning: recipe step ingredients
    ('perm_create_recipe_step_ingredients', 'create.recipe_step_ingredients', 'Create recipe step ingredients'),
    ('perm_read_recipe_step_ingredients', 'read.recipe_step_ingredients', 'Read recipe step ingredients'),
    ('perm_search_recipe_step_ingredients', 'search.recipe_step_ingredients', 'Search recipe step ingredients'),
    ('perm_update_recipe_step_ingredients', 'update.recipe_step_ingredients', 'Update recipe step ingredients'),
    ('perm_archive_recipe_step_ingredients', 'archive.recipe_step_ingredients', 'Archive recipe step ingredients'),
    -- meal planning: recipe step completion conditions
    ('perm_create_recipe_step_completion_conditions', 'create.recipe_step_completion_conditions', 'Create recipe step completion conditions'),
    ('perm_read_recipe_step_completion_conditions', 'read.recipe_step_completion_conditions', 'Read recipe step completion conditions'),
    ('perm_search_recipe_step_completion_conditions', 'search.recipe_step_completion_conditions', 'Search recipe step completion conditions'),
    ('perm_update_recipe_step_completion_conditions', 'update.recipe_step_completion_conditions', 'Update recipe step completion conditions'),
    ('perm_archive_recipe_step_completion_conditions', 'archive.recipe_step_completion_conditions', 'Archive recipe step completion conditions'),
    -- meal planning: recipe step products
    ('perm_create_recipe_step_products', 'create.recipe_step_products', 'Create recipe step products'),
    ('perm_read_recipe_step_products', 'read.recipe_step_products', 'Read recipe step products'),
    ('perm_search_recipe_step_products', 'search.recipe_step_products', 'Search recipe step products'),
    ('perm_update_recipe_step_products', 'update.recipe_step_products', 'Update recipe step products'),
    ('perm_archive_recipe_step_products', 'archive.recipe_step_products', 'Archive recipe step products'),
    -- meal planning: meal plans
    ('perm_create_meal_plans', 'create.meal_plans', 'Create meal plans'),
    ('perm_read_meal_plans', 'read.meal_plans', 'Read meal plans'),
    ('perm_search_meal_plans', 'search.meal_plans', 'Search meal plans'),
    ('perm_update_meal_plans', 'update.meal_plans', 'Update meal plans'),
    ('perm_archive_meal_plans', 'archive.meal_plans', 'Archive meal plans'),
    -- meal planning: meal plan events
    ('perm_create_meal_plan_events', 'create.meal_plan_events', 'Create meal plan events'),
    ('perm_read_meal_plan_events', 'read.meal_plan_events', 'Read meal plan events'),
    ('perm_update_meal_plan_events', 'update.meal_plan_events', 'Update meal plan events'),
    ('perm_archive_meal_plan_events', 'archive.meal_plan_events', 'Archive meal plan events'),
    -- meal planning: meal plan options
    ('perm_create_meal_plan_options', 'create.meal_plan_options', 'Create meal plan options'),
    ('perm_read_meal_plan_options', 'read.meal_plan_options', 'Read meal plan options'),
    ('perm_search_meal_plan_options', 'search.meal_plan_options', 'Search meal plan options'),
    ('perm_update_meal_plan_options', 'update.meal_plan_options', 'Update meal plan options'),
    ('perm_archive_meal_plan_options', 'archive.meal_plan_options', 'Archive meal plan options'),
    -- meal planning: meal plan tasks
    ('perm_create_meal_plan_tasks', 'create.meal_plan_tasks', 'Create meal plan tasks'),
    ('perm_read_meal_plan_tasks', 'read.meal_plan_tasks', 'Read meal plan tasks'),
    ('perm_update_meal_plan_tasks', 'update.meal_plan_tasks', 'Update meal plan tasks'),
    -- meal planning: meal plan grocery list items
    ('perm_create_meal_plan_grocery_list_items', 'create.meal_plan_grocery_list_items', 'Create meal plan grocery list items'),
    ('perm_read_meal_plan_grocery_list_items', 'read.meal_plan_grocery_list_items', 'Read meal plan grocery list items'),
    ('perm_update_meal_plan_grocery_list_items', 'update.meal_plan_grocery_list_items', 'Update meal plan grocery list items'),
    ('perm_archive_meal_plan_grocery_list_items', 'archive.meal_plan_grocery_list_items', 'Archive meal plan grocery list items'),
    -- meal planning: meal plan option votes
    ('perm_create_meal_plan_option_votes', 'create.meal_plan_option_votes', 'Create meal plan option votes'),
    ('perm_read_meal_plan_option_votes', 'read.meal_plan_option_votes', 'Read meal plan option votes'),
    ('perm_search_meal_plan_option_votes', 'search.meal_plan_option_votes', 'Search meal plan option votes'),
    ('perm_update_meal_plan_option_votes', 'update.meal_plan_option_votes', 'Update meal plan option votes'),
    ('perm_archive_meal_plan_option_votes', 'archive.meal_plan_option_votes', 'Archive meal plan option votes'),
    -- meal planning: meal plan recipe option selections
    ('perm_create_meal_plan_recipe_option_selections', 'create.meal_plan_recipe_option_selections', 'Create meal plan recipe option selections'),
    ('perm_read_meal_plan_recipe_option_selections', 'read.meal_plan_recipe_option_selections', 'Read meal plan recipe option selections'),
    ('perm_update_meal_plan_recipe_option_selections', 'update.meal_plan_recipe_option_selections', 'Update meal plan recipe option selections'),
    ('perm_archive_meal_plan_recipe_option_selections', 'archive.meal_plan_recipe_option_selections', 'Archive meal plan recipe option selections'),
    -- meal planning: user ingredient preferences
    ('perm_create_user_ingredient_preferences', 'create.user_ingredient_preferences', 'Create user ingredient preferences'),
    ('perm_read_user_ingredient_preferences', 'read.user_ingredient_preferences', 'Read user ingredient preferences'),
    ('perm_update_user_ingredient_preferences', 'update.user_ingredient_preferences', 'Update user ingredient preferences'),
    ('perm_archive_user_ingredient_preferences', 'archive.user_ingredient_preferences', 'Archive user ingredient preferences'),
    -- meal planning: account instrument ownerships
    ('perm_create_account_instrument_ownerships', 'create.account_instrument_ownerships', 'Create account instrument ownerships'),
    ('perm_read_account_instrument_ownerships', 'read.account_instrument_ownerships', 'Read account instrument ownerships'),
    ('perm_update_account_instrument_ownerships', 'update.account_instrument_ownerships', 'Update account instrument ownerships'),
    ('perm_archive_account_instrument_ownerships', 'archive.account_instrument_ownerships', 'Archive account instrument ownerships'),
    -- meal planning: recipe ratings
    ('perm_create_recipe_ratings', 'create.recipe_ratings', 'Create recipe ratings'),
    ('perm_read_recipe_ratings', 'read.recipe_ratings', 'Read recipe ratings'),
    ('perm_update_recipe_ratings', 'update.recipe_ratings', 'Update recipe ratings'),
    ('perm_archive_recipe_ratings', 'archive.recipe_ratings', 'Archive recipe ratings'),
    -- meal planning: meal lists
    ('perm_create_meal_lists', 'create.meal_lists', 'Create meal lists'),
    ('perm_read_meal_lists', 'read.meal_lists', 'Read meal lists'),
    ('perm_update_meal_lists', 'update.meal_lists', 'Update meal lists'),
    ('perm_archive_meal_lists', 'archive.meal_lists', 'Archive meal lists'),
    -- meal planning: recipe lists
    ('perm_create_recipe_lists', 'create.recipe_lists', 'Create recipe lists'),
    ('perm_read_recipe_lists', 'read.recipe_lists', 'Read recipe lists'),
    ('perm_update_recipe_lists', 'update.recipe_lists', 'Update recipe lists'),
    ('perm_archive_recipe_lists', 'archive.recipe_lists', 'Archive recipe lists');

-- =============================================================================
-- SEED DATA: role-permission mappings (direct assignments only)
-- Inheritance is resolved at query time via user_role_hierarchy
-- =============================================================================

-- service_admin direct permissions (27 + 42 from data admin = 69)
-- ServiceAdminPermissions
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
    ('urp_sa_13', 'role_service_admin', 'perm_update_recipe_status'),
    ('urp_sa_14', 'role_service_admin', 'perm_create_service_settings'),
    ('urp_sa_15', 'role_service_admin', 'perm_create_meal_plan_tasks'),
    ('urp_sa_16', 'role_service_admin', 'perm_create_meal_plan_grocery_list_items'),
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
    ('urp_sa_27', 'role_service_admin', 'perm_archive_subscriptions'),
    -- ServiceDataAdminPermissions also assigned to service_admin (per old init code)
    ('urp_sa_28', 'role_service_admin', 'perm_create_valid_instruments'),
    ('urp_sa_29', 'role_service_admin', 'perm_update_valid_instruments'),
    ('urp_sa_30', 'role_service_admin', 'perm_archive_valid_instruments'),
    ('urp_sa_31', 'role_service_admin', 'perm_create_valid_vessels'),
    ('urp_sa_32', 'role_service_admin', 'perm_update_valid_vessels'),
    ('urp_sa_33', 'role_service_admin', 'perm_archive_valid_vessels'),
    ('urp_sa_34', 'role_service_admin', 'perm_create_valid_ingredients'),
    ('urp_sa_35', 'role_service_admin', 'perm_update_valid_ingredients'),
    ('urp_sa_36', 'role_service_admin', 'perm_archive_valid_ingredients'),
    ('urp_sa_37', 'role_service_admin', 'perm_create_valid_ingredient_groups'),
    ('urp_sa_38', 'role_service_admin', 'perm_update_valid_ingredient_groups'),
    ('urp_sa_39', 'role_service_admin', 'perm_archive_valid_ingredient_groups'),
    ('urp_sa_40', 'role_service_admin', 'perm_create_valid_preparations'),
    ('urp_sa_41', 'role_service_admin', 'perm_update_valid_preparations'),
    ('urp_sa_42', 'role_service_admin', 'perm_archive_valid_preparations'),
    ('urp_sa_43', 'role_service_admin', 'perm_create_measurement_units'),
    ('urp_sa_44', 'role_service_admin', 'perm_update_measurement_units'),
    ('urp_sa_45', 'role_service_admin', 'perm_archive_measurement_units'),
    ('urp_sa_46', 'role_service_admin', 'perm_create_measurement_conversions'),
    ('urp_sa_47', 'role_service_admin', 'perm_update_measurement_conversions'),
    ('urp_sa_48', 'role_service_admin', 'perm_archive_measurement_conversions'),
    ('urp_sa_49', 'role_service_admin', 'perm_create_valid_ingredient_preparations'),
    ('urp_sa_50', 'role_service_admin', 'perm_update_valid_ingredient_preparations'),
    ('urp_sa_51', 'role_service_admin', 'perm_archive_valid_ingredient_preparations'),
    ('urp_sa_52', 'role_service_admin', 'perm_create_valid_prep_task_configs'),
    ('urp_sa_53', 'role_service_admin', 'perm_update_valid_prep_task_configs'),
    ('urp_sa_54', 'role_service_admin', 'perm_archive_valid_prep_task_configs'),
    ('urp_sa_55', 'role_service_admin', 'perm_create_valid_ingredient_state_ingredients'),
    ('urp_sa_56', 'role_service_admin', 'perm_update_valid_ingredient_state_ingredients'),
    ('urp_sa_57', 'role_service_admin', 'perm_archive_valid_ingredient_state_ingredients'),
    ('urp_sa_58', 'role_service_admin', 'perm_create_valid_preparation_instruments'),
    ('urp_sa_59', 'role_service_admin', 'perm_update_valid_preparation_instruments'),
    ('urp_sa_60', 'role_service_admin', 'perm_archive_valid_preparation_instruments'),
    ('urp_sa_61', 'role_service_admin', 'perm_create_valid_preparation_vessels'),
    ('urp_sa_62', 'role_service_admin', 'perm_update_valid_preparation_vessels'),
    ('urp_sa_63', 'role_service_admin', 'perm_archive_valid_preparation_vessels'),
    ('urp_sa_64', 'role_service_admin', 'perm_create_valid_ingredient_measurement_units'),
    ('urp_sa_65', 'role_service_admin', 'perm_update_valid_ingredient_measurement_units'),
    ('urp_sa_66', 'role_service_admin', 'perm_archive_valid_ingredient_measurement_units'),
    ('urp_sa_67', 'role_service_admin', 'perm_create_valid_ingredient_states'),
    ('urp_sa_68', 'role_service_admin', 'perm_update_valid_ingredient_states'),
    ('urp_sa_69', 'role_service_admin', 'perm_archive_valid_ingredient_states');

-- service_data_admin direct permissions (42)
INSERT INTO user_role_permissions (id, role_id, permission_id) VALUES
    ('urp_sda_1', 'role_service_data_admin', 'perm_create_valid_instruments'),
    ('urp_sda_2', 'role_service_data_admin', 'perm_update_valid_instruments'),
    ('urp_sda_3', 'role_service_data_admin', 'perm_archive_valid_instruments'),
    ('urp_sda_4', 'role_service_data_admin', 'perm_create_valid_vessels'),
    ('urp_sda_5', 'role_service_data_admin', 'perm_update_valid_vessels'),
    ('urp_sda_6', 'role_service_data_admin', 'perm_archive_valid_vessels'),
    ('urp_sda_7', 'role_service_data_admin', 'perm_create_valid_ingredients'),
    ('urp_sda_8', 'role_service_data_admin', 'perm_update_valid_ingredients'),
    ('urp_sda_9', 'role_service_data_admin', 'perm_archive_valid_ingredients'),
    ('urp_sda_10', 'role_service_data_admin', 'perm_create_valid_ingredient_groups'),
    ('urp_sda_11', 'role_service_data_admin', 'perm_update_valid_ingredient_groups'),
    ('urp_sda_12', 'role_service_data_admin', 'perm_archive_valid_ingredient_groups'),
    ('urp_sda_13', 'role_service_data_admin', 'perm_create_valid_preparations'),
    ('urp_sda_14', 'role_service_data_admin', 'perm_update_valid_preparations'),
    ('urp_sda_15', 'role_service_data_admin', 'perm_archive_valid_preparations'),
    ('urp_sda_16', 'role_service_data_admin', 'perm_create_measurement_units'),
    ('urp_sda_17', 'role_service_data_admin', 'perm_update_measurement_units'),
    ('urp_sda_18', 'role_service_data_admin', 'perm_archive_measurement_units'),
    ('urp_sda_19', 'role_service_data_admin', 'perm_create_measurement_conversions'),
    ('urp_sda_20', 'role_service_data_admin', 'perm_update_measurement_conversions'),
    ('urp_sda_21', 'role_service_data_admin', 'perm_archive_measurement_conversions'),
    ('urp_sda_22', 'role_service_data_admin', 'perm_create_valid_ingredient_preparations'),
    ('urp_sda_23', 'role_service_data_admin', 'perm_update_valid_ingredient_preparations'),
    ('urp_sda_24', 'role_service_data_admin', 'perm_archive_valid_ingredient_preparations'),
    ('urp_sda_25', 'role_service_data_admin', 'perm_create_valid_prep_task_configs'),
    ('urp_sda_26', 'role_service_data_admin', 'perm_update_valid_prep_task_configs'),
    ('urp_sda_27', 'role_service_data_admin', 'perm_archive_valid_prep_task_configs'),
    ('urp_sda_28', 'role_service_data_admin', 'perm_create_valid_ingredient_state_ingredients'),
    ('urp_sda_29', 'role_service_data_admin', 'perm_update_valid_ingredient_state_ingredients'),
    ('urp_sda_30', 'role_service_data_admin', 'perm_archive_valid_ingredient_state_ingredients'),
    ('urp_sda_31', 'role_service_data_admin', 'perm_create_valid_preparation_instruments'),
    ('urp_sda_32', 'role_service_data_admin', 'perm_update_valid_preparation_instruments'),
    ('urp_sda_33', 'role_service_data_admin', 'perm_archive_valid_preparation_instruments'),
    ('urp_sda_34', 'role_service_data_admin', 'perm_create_valid_preparation_vessels'),
    ('urp_sda_35', 'role_service_data_admin', 'perm_update_valid_preparation_vessels'),
    ('urp_sda_36', 'role_service_data_admin', 'perm_archive_valid_preparation_vessels'),
    ('urp_sda_37', 'role_service_data_admin', 'perm_create_valid_ingredient_measurement_units'),
    ('urp_sda_38', 'role_service_data_admin', 'perm_update_valid_ingredient_measurement_units'),
    ('urp_sda_39', 'role_service_data_admin', 'perm_archive_valid_ingredient_measurement_units'),
    ('urp_sda_40', 'role_service_data_admin', 'perm_create_valid_ingredient_states'),
    ('urp_sda_41', 'role_service_data_admin', 'perm_update_valid_ingredient_states'),
    ('urp_sda_42', 'role_service_data_admin', 'perm_archive_valid_ingredient_states');

-- account_admin direct permissions (47)
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
    ('urp_aa_13', 'role_account_admin', 'perm_create_meal_plans'),
    ('urp_aa_14', 'role_account_admin', 'perm_update_meal_plans'),
    ('urp_aa_15', 'role_account_admin', 'perm_archive_meal_plans'),
    ('urp_aa_16', 'role_account_admin', 'perm_create_meal_plan_events'),
    ('urp_aa_17', 'role_account_admin', 'perm_update_meal_plan_events'),
    ('urp_aa_18', 'role_account_admin', 'perm_archive_meal_plan_events'),
    ('urp_aa_19', 'role_account_admin', 'perm_create_meal_plan_options'),
    ('urp_aa_20', 'role_account_admin', 'perm_update_meal_plan_options'),
    ('urp_aa_21', 'role_account_admin', 'perm_archive_meal_plan_options'),
    ('urp_aa_22', 'role_account_admin', 'perm_create_account_instrument_ownerships'),
    ('urp_aa_23', 'role_account_admin', 'perm_update_account_instrument_ownerships'),
    ('urp_aa_24', 'role_account_admin', 'perm_archive_account_instrument_ownerships'),
    ('urp_aa_25', 'role_account_admin', 'perm_create_webhook_trigger_configs'),
    ('urp_aa_26', 'role_account_admin', 'perm_archive_webhook_trigger_configs'),
    ('urp_aa_27', 'role_account_admin', 'perm_create_webhook_trigger_events'),
    ('urp_aa_28', 'role_account_admin', 'perm_read_webhook_trigger_events'),
    ('urp_aa_29', 'role_account_admin', 'perm_update_webhook_trigger_events'),
    ('urp_aa_30', 'role_account_admin', 'perm_archive_webhook_trigger_events'),
    ('urp_aa_31', 'role_account_admin', 'perm_create_meal_lists'),
    ('urp_aa_32', 'role_account_admin', 'perm_read_meal_lists'),
    ('urp_aa_33', 'role_account_admin', 'perm_update_meal_lists'),
    ('urp_aa_34', 'role_account_admin', 'perm_archive_meal_lists'),
    ('urp_aa_35', 'role_account_admin', 'perm_create_recipe_lists'),
    ('urp_aa_36', 'role_account_admin', 'perm_read_recipe_lists'),
    ('urp_aa_37', 'role_account_admin', 'perm_update_recipe_lists'),
    ('urp_aa_38', 'role_account_admin', 'perm_archive_recipe_lists'),
    ('urp_aa_39', 'role_account_admin', 'perm_create_checkout_sessions'),
    ('urp_aa_40', 'role_account_admin', 'perm_cancel_subscriptions'),
    ('urp_aa_41', 'role_account_admin', 'perm_read_purchases'),
    ('urp_aa_42', 'role_account_admin', 'perm_read_payment_history'),
    ('urp_aa_43', 'role_account_admin', 'perm_read_subscriptions');

-- account_member direct permissions
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
    ('urp_am_12', 'role_account_member', 'perm_create_meals'),
    ('urp_am_13', 'role_account_member', 'perm_read_meals'),
    ('urp_am_14', 'role_account_member', 'perm_update_meals'),
    ('urp_am_15', 'role_account_member', 'perm_archive_meals'),
    ('urp_am_16', 'role_account_member', 'perm_create_recipes'),
    ('urp_am_17', 'role_account_member', 'perm_read_recipes'),
    ('urp_am_18', 'role_account_member', 'perm_search_recipes'),
    ('urp_am_19', 'role_account_member', 'perm_update_recipes'),
    ('urp_am_20', 'role_account_member', 'perm_archive_recipes'),
    ('urp_am_21', 'role_account_member', 'perm_create_recipe_steps'),
    ('urp_am_22', 'role_account_member', 'perm_read_recipe_steps'),
    ('urp_am_23', 'role_account_member', 'perm_search_recipe_steps'),
    ('urp_am_24', 'role_account_member', 'perm_update_recipe_steps'),
    ('urp_am_25', 'role_account_member', 'perm_archive_recipe_steps'),
    ('urp_am_26', 'role_account_member', 'perm_create_recipe_prep_tasks'),
    ('urp_am_27', 'role_account_member', 'perm_read_recipe_prep_tasks'),
    ('urp_am_28', 'role_account_member', 'perm_update_recipe_prep_tasks'),
    ('urp_am_29', 'role_account_member', 'perm_archive_recipe_prep_tasks'),
    ('urp_am_30', 'role_account_member', 'perm_create_recipe_step_instruments'),
    ('urp_am_31', 'role_account_member', 'perm_read_recipe_step_instruments'),
    ('urp_am_32', 'role_account_member', 'perm_search_recipe_step_instruments'),
    ('urp_am_33', 'role_account_member', 'perm_update_recipe_step_instruments'),
    ('urp_am_34', 'role_account_member', 'perm_archive_recipe_step_instruments'),
    ('urp_am_35', 'role_account_member', 'perm_create_recipe_step_vessels'),
    ('urp_am_36', 'role_account_member', 'perm_read_recipe_step_vessels'),
    ('urp_am_37', 'role_account_member', 'perm_search_recipe_step_vessels'),
    ('urp_am_38', 'role_account_member', 'perm_update_recipe_step_vessels'),
    ('urp_am_39', 'role_account_member', 'perm_archive_recipe_step_vessels'),
    ('urp_am_40', 'role_account_member', 'perm_create_recipe_step_ingredients'),
    ('urp_am_41', 'role_account_member', 'perm_read_recipe_step_ingredients'),
    ('urp_am_42', 'role_account_member', 'perm_search_recipe_step_ingredients'),
    ('urp_am_43', 'role_account_member', 'perm_update_recipe_step_ingredients'),
    ('urp_am_44', 'role_account_member', 'perm_archive_recipe_step_ingredients'),
    ('urp_am_45', 'role_account_member', 'perm_create_recipe_step_completion_conditions'),
    ('urp_am_46', 'role_account_member', 'perm_read_recipe_step_completion_conditions'),
    ('urp_am_47', 'role_account_member', 'perm_search_recipe_step_completion_conditions'),
    ('urp_am_48', 'role_account_member', 'perm_update_recipe_step_completion_conditions'),
    ('urp_am_49', 'role_account_member', 'perm_archive_recipe_step_completion_conditions'),
    ('urp_am_50', 'role_account_member', 'perm_create_recipe_step_products'),
    ('urp_am_51', 'role_account_member', 'perm_read_recipe_step_products'),
    ('urp_am_52', 'role_account_member', 'perm_search_recipe_step_products'),
    ('urp_am_53', 'role_account_member', 'perm_update_recipe_step_products'),
    ('urp_am_54', 'role_account_member', 'perm_archive_recipe_step_products'),
    ('urp_am_55', 'role_account_member', 'perm_read_valid_instruments'),
    ('urp_am_56', 'role_account_member', 'perm_search_valid_instruments'),
    ('urp_am_57', 'role_account_member', 'perm_read_valid_vessels'),
    ('urp_am_58', 'role_account_member', 'perm_search_valid_vessels'),
    ('urp_am_59', 'role_account_member', 'perm_read_valid_ingredients'),
    ('urp_am_60', 'role_account_member', 'perm_search_valid_ingredients'),
    ('urp_am_61', 'role_account_member', 'perm_read_valid_ingredient_groups'),
    ('urp_am_62', 'role_account_member', 'perm_search_valid_ingredient_groups'),
    ('urp_am_63', 'role_account_member', 'perm_read_valid_preparations'),
    ('urp_am_64', 'role_account_member', 'perm_search_valid_preparations'),
    ('urp_am_65', 'role_account_member', 'perm_read_measurement_units'),
    ('urp_am_66', 'role_account_member', 'perm_search_measurement_units'),
    ('urp_am_67', 'role_account_member', 'perm_read_measurement_conversions'),
    ('urp_am_68', 'role_account_member', 'perm_read_valid_ingredient_preparations'),
    ('urp_am_69', 'role_account_member', 'perm_search_valid_ingredient_preparations'),
    ('urp_am_70', 'role_account_member', 'perm_read_valid_ingredient_state_ingredients'),
    ('urp_am_71', 'role_account_member', 'perm_search_valid_ingredient_state_ingredients'),
    ('urp_am_72', 'role_account_member', 'perm_read_valid_preparation_instruments'),
    ('urp_am_73', 'role_account_member', 'perm_search_valid_preparation_instruments'),
    ('urp_am_74', 'role_account_member', 'perm_read_valid_preparation_vessels'),
    ('urp_am_75', 'role_account_member', 'perm_search_valid_preparation_vessels'),
    ('urp_am_76', 'role_account_member', 'perm_read_valid_ingredient_measurement_units'),
    ('urp_am_77', 'role_account_member', 'perm_search_valid_ingredient_measurement_units'),
    ('urp_am_78', 'role_account_member', 'perm_read_meal_plans'),
    ('urp_am_79', 'role_account_member', 'perm_search_meal_plans'),
    ('urp_am_80', 'role_account_member', 'perm_read_meal_plan_events'),
    ('urp_am_81', 'role_account_member', 'perm_read_meal_plan_options'),
    ('urp_am_82', 'role_account_member', 'perm_search_meal_plan_options'),
    ('urp_am_83', 'role_account_member', 'perm_read_valid_ingredient_states'),
    ('urp_am_84', 'role_account_member', 'perm_read_meal_plan_grocery_list_items'),
    ('urp_am_85', 'role_account_member', 'perm_update_meal_plan_grocery_list_items'),
    ('urp_am_86', 'role_account_member', 'perm_archive_meal_plan_grocery_list_items'),
    ('urp_am_87', 'role_account_member', 'perm_create_meal_plan_option_votes'),
    ('urp_am_88', 'role_account_member', 'perm_read_meal_plan_option_votes'),
    ('urp_am_89', 'role_account_member', 'perm_search_meal_plan_option_votes'),
    ('urp_am_90', 'role_account_member', 'perm_update_meal_plan_option_votes'),
    ('urp_am_91', 'role_account_member', 'perm_archive_meal_plan_option_votes'),
    ('urp_am_92', 'role_account_member', 'perm_create_meal_plan_recipe_option_selections'),
    ('urp_am_93', 'role_account_member', 'perm_read_meal_plan_recipe_option_selections'),
    ('urp_am_94', 'role_account_member', 'perm_update_meal_plan_recipe_option_selections'),
    ('urp_am_95', 'role_account_member', 'perm_archive_meal_plan_recipe_option_selections'),
    ('urp_am_96', 'role_account_member', 'perm_create_service_setting_configurations'),
    ('urp_am_97', 'role_account_member', 'perm_read_service_setting_configurations'),
    ('urp_am_98', 'role_account_member', 'perm_update_service_setting_configurations'),
    ('urp_am_99', 'role_account_member', 'perm_archive_service_setting_configurations'),
    ('urp_am_100', 'role_account_member', 'perm_read_meal_plan_tasks'),
    ('urp_am_101', 'role_account_member', 'perm_update_meal_plan_tasks'),
    ('urp_am_102', 'role_account_member', 'perm_create_user_ingredient_preferences'),
    ('urp_am_103', 'role_account_member', 'perm_read_user_ingredient_preferences'),
    ('urp_am_104', 'role_account_member', 'perm_update_user_ingredient_preferences'),
    ('urp_am_105', 'role_account_member', 'perm_archive_user_ingredient_preferences'),
    ('urp_am_106', 'role_account_member', 'perm_read_account_instrument_ownerships'),
    ('urp_am_107', 'role_account_member', 'perm_create_recipe_ratings'),
    ('urp_am_108', 'role_account_member', 'perm_read_recipe_ratings'),
    ('urp_am_109', 'role_account_member', 'perm_create_comments'),
    ('urp_am_110', 'role_account_member', 'perm_read_comments'),
    ('urp_am_111', 'role_account_member', 'perm_update_comments'),
    ('urp_am_112', 'role_account_member', 'perm_archive_comments'),
    ('urp_am_113', 'role_account_member', 'perm_update_recipe_ratings'),
    ('urp_am_114', 'role_account_member', 'perm_archive_recipe_ratings'),
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
    ('urp_am_125', 'role_account_member', 'perm_read_valid_prep_task_configs'),
    ('urp_am_126', 'role_account_member', 'perm_create_checkout_sessions'),
    ('urp_am_127', 'role_account_member', 'perm_cancel_subscriptions'),
    ('urp_am_128', 'role_account_member', 'perm_read_purchases'),
    ('urp_am_129', 'role_account_member', 'perm_read_payment_history'),
    ('urp_am_130', 'role_account_member', 'perm_read_subscriptions');
