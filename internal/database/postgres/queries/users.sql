-- name: GetUserByID :one
	SELECT
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.hashed_password,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.user_account_status,
		users.user_account_status_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on
	FROM users
	WHERE users.archived_on IS NULL
	AND users.id = $1;

-- name: GetUserWithVerified2FA :one
	SELECT
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.hashed_password,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.user_account_status,
		users.user_account_status_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on
	FROM users
	WHERE users.archived_on IS NULL
	AND users.id = $1
	AND users.two_factor_secret_verified_on IS NOT NULL;

-- name: UserHasStatus :one
SELECT EXISTS ( SELECT users.id FROM users WHERE users.archived_on IS NULL AND users.id = $1 AND (users.user_account_status = $2 OR users.user_account_status = $3) );

-- name: GetUserByUsername :one
SELECT
    users.id,
    users.username,
    users.email_address,
    users.avatar_src,
    users.hashed_password,
    users.requires_password_change,
    users.password_last_changed_on,
    users.two_factor_secret,
    users.two_factor_secret_verified_on,
    users.service_roles,
    users.user_account_status,
    users.user_account_status_explanation,
    users.birth_day,
    users.birth_month,
    users.created_on,
    users.last_updated_on,
    users.archived_on
FROM users
WHERE users.archived_on IS NULL
AND users.username = $1;

-- name: GetAdminUserByUsername :one
SELECT
    users.id,
    users.username,
    users.email_address,
    users.avatar_src,
    users.hashed_password,
    users.requires_password_change,
    users.password_last_changed_on,
    users.two_factor_secret,
    users.two_factor_secret_verified_on,
    users.service_roles,
    users.user_account_status,
    users.user_account_status_explanation,
    users.birth_day,
    users.birth_month,
    users.created_on,
    users.last_updated_on,
    users.archived_on
FROM users
WHERE users.archived_on IS NULL
AND users.service_roles ILIKE '%service_admin%'
AND users.username = $1
AND users.two_factor_secret_verified_on IS NOT NULL;

-- name: GetUserIDByEmail :one
SELECT
    users.id,
    users.username,
    users.email_address,
    users.avatar_src,
    users.hashed_password,
    users.requires_password_change,
    users.password_last_changed_on,
    users.two_factor_secret,
    users.two_factor_secret_verified_on,
    users.service_roles,
    users.user_account_status,
    users.user_account_status_explanation,
    users.birth_day,
    users.birth_month,
    users.created_on,
    users.last_updated_on,
    users.archived_on
FROM users
WHERE users.archived_on IS NULL
AND users.email_address = $1;

-- name: SearchForUserByUsername :many
SELECT
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_on,
	users.two_factor_secret,
	users.two_factor_secret_verified_on,
	users.service_roles,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_on,
	users.last_updated_on,
	users.archived_on
FROM users
WHERE users.username ILIKE $1
AND users.archived_on IS NULL
AND users.two_factor_secret_verified_on IS NOT NULL;

-- name: GetAllUsersCount :one
SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL;

-- name: CreateUser :exec
INSERT INTO users (id,username,email_address,hashed_password,two_factor_secret,avatar_src,user_account_status,birth_day,birth_month,service_roles) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: CreateHouseholdMembershipForNewUser :exec
INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,default_household,household_roles)
	VALUES ($1,$2,$3,$4,$5);

-- name: UpdateUser :exec
UPDATE users SET
	username = $1,
	hashed_password = $2,
	avatar_src = $3,
	two_factor_secret = $4,
	two_factor_secret_verified_on = $5,
	birth_day = $6,
	birth_month = $7,
	last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $8;

-- name: UpdateUserPassword :exec
UPDATE users SET
	hashed_password = $1,
	requires_password_change = $2,
	password_last_changed_on = extract(epoch FROM NOW()),
	last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $3;

-- name: UpdateUserTwoFactorSecret :exec
UPDATE users SET
	two_factor_secret_verified_on = $1,
	two_factor_secret = $2
WHERE archived_on IS NULL
AND id = $3;

-- name: MarkUserTwoFactorSecretAsVerified :exec
UPDATE users SET
	two_factor_secret_verified_on = extract(epoch FROM NOW()),
	user_account_status = $1
WHERE archived_on IS NULL
AND id = $2;

-- name: ArchiveUser :exec
UPDATE users SET
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $1;

-- name: ArchiveMemberships :exec
UPDATE household_user_memberships SET
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND belongs_to_user = $1;
