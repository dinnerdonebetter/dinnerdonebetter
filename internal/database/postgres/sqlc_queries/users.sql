-- name: AcceptPrivacyPolicyForUser :exec

UPDATE users SET last_accepted_privacy_policy = NOW() WHERE archived_at IS NULL AND id = sqlc.arg(id);

-- name: AcceptTermsOfServiceForUser :exec

UPDATE users SET last_accepted_terms_of_service = NOW() WHERE archived_at IS NULL AND id = sqlc.arg(id);

-- name: ArchiveUser :execrows

UPDATE users SET archived_at = NOW() WHERE archived_at IS NULL AND id = sqlc.arg(id);

-- name: ArchiveUserMemberships :execrows

UPDATE household_user_memberships SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_user = sqlc.arg(id);

-- name: CreateUser :exec

INSERT
INTO
	users
		(
			id,
			username,
			avatar_src,
			email_address,
			hashed_password,
			requires_password_change,
			two_factor_secret,
			two_factor_secret_verified_at,
			service_role,
			user_account_status,
			user_account_status_explanation,
			birthday,
			email_address_verification_token,
			first_name,
			last_name
		)
VALUES
	(
		sqlc.arg(id),
		sqlc.arg(username),
		sqlc.arg(avatar_src),
		sqlc.arg(email_address),
		sqlc.arg(hashed_password),
		sqlc.arg(requires_password_change),
		sqlc.arg(two_factor_secret),
		sqlc.arg(two_factor_secret_verified_at),
		sqlc.arg(service_role),
		sqlc.arg(user_account_status),
		sqlc.arg(user_account_status_explanation),
		sqlc.arg(birthday),
		sqlc.arg(email_address_verification_token),
		sqlc.arg(first_name),
		sqlc.arg(last_name)
	);

-- name: GetAdminUserByUsername :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.service_role = 'service_admin'
	AND users.username = $1
	AND users.two_factor_secret_verified_at IS NOT NULL;

-- name: GetUserByEmail :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.email_address = $1;

-- name: GetUserByEmailAddressVerificationToken :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.email_address_verification_token = $1;

-- name: GetUserByID :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.id = $1;

-- name: GetUserByUsername :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.username = $1;

-- name: GetEmailVerificationTokenByUserID :one

SELECT
	users.email_address_verification_token
FROM users
WHERE users.archived_at IS NULL
    AND users.email_address_verified_at IS NULL
	AND users.id = $1;

-- name: GetUsers :many

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at,
    (
        SELECT
            COUNT(users.id)
        FROM
            users
        WHERE
            users.archived_at IS NULL
          AND users.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
          AND users.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
          AND (
                users.last_updated_at IS NULL
                OR users.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
          AND (
                users.last_updated_at IS NULL
                OR users.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) as filtered_count,
    (
        SELECT
            COUNT(users.id)
        FROM
            users
        WHERE
            users.archived_at IS NULL
    ) as total_count
FROM users
WHERE
    users.archived_at IS NULL
  AND users.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
  AND users.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
  AND (
        users.last_updated_at IS NULL
        OR users.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
    )
  AND (
        users.last_updated_at IS NULL
        OR users.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
    )
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetUserIDsNeedingIndexing :many

SELECT users.id
  FROM users
 WHERE (users.archived_at IS NULL)
       AND (
			(
				users.last_indexed_at IS NULL
			)
			OR users.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetUserWithUnverifiedTwoFactor :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.id = $1
	AND users.two_factor_secret_verified_at IS NULL;

-- name: GetUserWithVerifiedTwoFactor :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.id = $1
	AND users.two_factor_secret_verified_at IS NOT NULL;

-- name: MarkEmailAddressAsVerified :exec

UPDATE users SET
	email_address_verified_at = NOW(),
	last_updated_at = NOW()
WHERE email_address_verified_at IS NULL
	AND id = $1
	AND email_address_verification_token = $2;

-- name: MarkTwoFactorSecretAsUnverified :exec

UPDATE users SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;

-- name: MarkTwoFactorSecretAsVerified :exec

UPDATE users SET
	two_factor_secret_verified_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;

-- name: SearchUsersByUsername :many

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.username ILIKE '%' || sqlc.arg(username)::text || '%'
AND users.archived_at IS NULL;

-- name: UpdateUser :execrows

UPDATE users SET
	username = $1,
	first_name = $2,
	last_name = $3,
	hashed_password = $4,
	avatar_src = $5,
	birthday = $6,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $7;

-- name: UpdateUserAvatarSrc :execrows

UPDATE users SET
	avatar_src = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;

-- name: UpdateUserDetails :execrows

UPDATE users SET
	first_name = $1,
	last_name = $2,
	birthday = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $4;

-- name: UpdateUserEmailAddress :execrows

UPDATE users SET
	email_address = $1,
	email_address_verified_at = NULL,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;

-- name: UpdateUserLastIndexedAt :execrows

UPDATE users SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;

-- name: UpdateUserPassword :execrows

UPDATE users SET
	hashed_password = $1,
	password_last_changed_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;

-- name: UpdateUserTwoFactorSecret :execrows

UPDATE users SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;

-- name: UpdateUserUsername :execrows

UPDATE users SET
	username = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
