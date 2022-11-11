SELECT
	users.id,
	users.username,
	users.email_address,
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
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.id = $1
	AND users.two_factor_secret_verified_at IS NOT NULL;