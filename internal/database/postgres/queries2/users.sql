-- name: AcceptPrivacyPolicyForUser :exec

UPDATE users SET
	last_accepted_privacy_policy = NOW()
WHERE archived_at IS NULL
	AND id = $1;


-- name: AcceptTermsOfServiceForUser :exec

UPDATE users SET
	last_accepted_terms_of_service = NOW()
WHERE archived_at IS NULL
	AND id = $1;


-- name: ArchiveUser :exec

UPDATE users SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;


-- name: ArchiveUserMemberships :exec

UPDATE household_user_memberships SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $1;


-- name: CreateUser :exec

INSERT INTO users (id,first_name,last_name,username,email_address,hashed_password,two_factor_secret,avatar_src,user_account_status,birthday,service_role,email_address_verification_token) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);


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


-- name: UpdateUser :exec

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


-- name: UpdateUserAvatarSrc :exec

UPDATE users SET
	avatar_src = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;


-- name: UpdateUserDetails :exec

UPDATE users SET
	first_name = $1,
	last_name = $2,
	birthday = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $4;


-- name: UpdateUserEmailAddress :exec

UPDATE users SET
	email_address = $1,
	email_address_verified_at = NULL,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;


-- name: UpdateUserLastIndexedAt :exec

UPDATE users SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;


-- name: UpdateUserPassword :exec

UPDATE users SET
	hashed_password = $1,
	password_last_changed_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;


-- name: UpdateUserTwoFactorSecret :exec

UPDATE users SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;


-- name: UpdateUserUsername :exec

UPDATE users SET
	username = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
