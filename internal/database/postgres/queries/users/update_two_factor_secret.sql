-- name: UpdateUserTwoFactorSecret :exec

UPDATE users SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
