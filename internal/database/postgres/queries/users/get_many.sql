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
