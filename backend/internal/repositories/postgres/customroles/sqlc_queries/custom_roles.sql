-- name: CreateCustomRole :exec
INSERT INTO custom_roles (
    id,
    name,
    description,
    scope,
    belongs_to_account,
    created_by
) VALUES (
    sqlc.arg(id),
    sqlc.arg(name),
    sqlc.arg(description),
    sqlc.arg(scope),
    sqlc.arg(belongs_to_account),
    sqlc.arg(created_by)
);

-- name: GetCustomRole :one
SELECT
    custom_roles.id,
    custom_roles.name,
    custom_roles.description,
    custom_roles.scope,
    custom_roles.belongs_to_account,
    custom_roles.created_by,
    custom_roles.created_at,
    custom_roles.last_updated_at,
    custom_roles.archived_at
FROM custom_roles
WHERE custom_roles.id = sqlc.arg(id)
    AND custom_roles.archived_at IS NULL;

-- name: ListServiceScopedCustomRoles :many
SELECT
    custom_roles.id,
    custom_roles.name,
    custom_roles.description,
    custom_roles.scope,
    custom_roles.belongs_to_account,
    custom_roles.created_by,
    custom_roles.created_at,
    custom_roles.last_updated_at,
    custom_roles.archived_at,
    (
        SELECT COUNT(custom_roles.id)
        FROM custom_roles
        WHERE custom_roles.scope = 'service'
            AND custom_roles.archived_at IS NULL
    ) AS filtered_count,
    (
        SELECT COUNT(custom_roles.id)
        FROM custom_roles
        WHERE custom_roles.scope = 'service'
            AND custom_roles.archived_at IS NULL
    ) AS total_count
FROM custom_roles
WHERE custom_roles.scope = 'service'
    AND custom_roles.archived_at IS NULL
ORDER BY custom_roles.created_at DESC
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: ListAccountScopedCustomRoles :many
SELECT
    custom_roles.id,
    custom_roles.name,
    custom_roles.description,
    custom_roles.scope,
    custom_roles.belongs_to_account,
    custom_roles.created_by,
    custom_roles.created_at,
    custom_roles.last_updated_at,
    custom_roles.archived_at,
    (
        SELECT COUNT(custom_roles.id)
        FROM custom_roles
        WHERE custom_roles.scope = 'account'
            AND custom_roles.belongs_to_account = sqlc.arg(belongs_to_account)
            AND custom_roles.archived_at IS NULL
    ) AS filtered_count,
    (
        SELECT COUNT(custom_roles.id)
        FROM custom_roles
        WHERE custom_roles.scope = 'account'
            AND custom_roles.belongs_to_account = sqlc.arg(belongs_to_account)
            AND custom_roles.archived_at IS NULL
    ) AS total_count
FROM custom_roles
WHERE custom_roles.scope = 'account'
    AND custom_roles.belongs_to_account = sqlc.arg(belongs_to_account)
    AND custom_roles.archived_at IS NULL
ORDER BY custom_roles.created_at DESC
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: UpdateCustomRole :exec
UPDATE custom_roles SET
    name = sqlc.arg(name),
    description = sqlc.arg(description),
    last_updated_at = NOW()
WHERE id = sqlc.arg(id)
    AND archived_at IS NULL;

-- name: ArchiveCustomRole :exec
UPDATE custom_roles SET
    archived_at = NOW()
WHERE id = sqlc.arg(id)
    AND archived_at IS NULL;

-- name: GetPermissionsForCustomRole :many
SELECT
    custom_role_permissions.id,
    custom_role_permissions.belongs_to_role,
    custom_role_permissions.permission,
    custom_role_permissions.created_at
FROM custom_role_permissions
WHERE custom_role_permissions.belongs_to_role = sqlc.arg(belongs_to_role);

-- name: CreateCustomRolePermission :exec
INSERT INTO custom_role_permissions (
    id,
    belongs_to_role,
    permission
) VALUES (
    sqlc.arg(id),
    sqlc.arg(belongs_to_role),
    sqlc.arg(permission)
);

-- name: GetAllCustomRolePermissions :many
SELECT
    custom_role_permissions.belongs_to_role,
    custom_role_permissions.permission
FROM custom_role_permissions
    JOIN custom_roles ON custom_roles.id = custom_role_permissions.belongs_to_role
WHERE custom_roles.archived_at IS NULL;

-- name: AssignCustomRoleToUser :exec
INSERT INTO user_custom_role_assignments (
    id,
    custom_role_id,
    user_id
) VALUES (
    sqlc.arg(id),
    sqlc.arg(custom_role_id),
    sqlc.arg(user_id)
);

-- name: AssignCustomRoleToMembership :exec
INSERT INTO user_custom_role_assignments (
    id,
    custom_role_id,
    account_membership_id
) VALUES (
    sqlc.arg(id),
    sqlc.arg(custom_role_id),
    sqlc.arg(account_membership_id)
);

-- name: UnassignCustomRole :exec
UPDATE user_custom_role_assignments SET
    archived_at = NOW()
WHERE id = sqlc.arg(id)
    AND archived_at IS NULL;

-- name: GetCustomRoleAssignmentsForUser :many
SELECT
    user_custom_role_assignments.id,
    user_custom_role_assignments.custom_role_id,
    custom_roles.name AS custom_role_name,
    user_custom_role_assignments.user_id,
    user_custom_role_assignments.account_membership_id,
    user_custom_role_assignments.created_at,
    user_custom_role_assignments.archived_at
FROM user_custom_role_assignments
    JOIN custom_roles ON custom_roles.id = user_custom_role_assignments.custom_role_id
WHERE user_custom_role_assignments.user_id = sqlc.arg(user_id)
    AND user_custom_role_assignments.archived_at IS NULL
    AND custom_roles.archived_at IS NULL;

-- name: GetCustomRoleAssignmentsForMembership :many
SELECT
    user_custom_role_assignments.id,
    user_custom_role_assignments.custom_role_id,
    custom_roles.name AS custom_role_name,
    user_custom_role_assignments.user_id,
    user_custom_role_assignments.account_membership_id,
    user_custom_role_assignments.created_at,
    user_custom_role_assignments.archived_at
FROM user_custom_role_assignments
    JOIN custom_roles ON custom_roles.id = user_custom_role_assignments.custom_role_id
WHERE user_custom_role_assignments.account_membership_id = sqlc.arg(account_membership_id)
    AND user_custom_role_assignments.archived_at IS NULL
    AND custom_roles.archived_at IS NULL;

-- name: GetServiceScopedRoleIDsForUser :many
SELECT user_custom_role_assignments.custom_role_id
FROM user_custom_role_assignments
    JOIN custom_roles ON custom_roles.id = user_custom_role_assignments.custom_role_id
WHERE user_custom_role_assignments.user_id = sqlc.arg(user_id)
    AND user_custom_role_assignments.archived_at IS NULL
    AND custom_roles.archived_at IS NULL
    AND custom_roles.scope = 'service';

-- name: GetAccountScopedRoleIDsForMembership :many
SELECT user_custom_role_assignments.custom_role_id
FROM user_custom_role_assignments
    JOIN custom_roles ON custom_roles.id = user_custom_role_assignments.custom_role_id
WHERE user_custom_role_assignments.account_membership_id = sqlc.arg(account_membership_id)
    AND user_custom_role_assignments.archived_at IS NULL
    AND custom_roles.archived_at IS NULL
    AND custom_roles.scope = 'account';
