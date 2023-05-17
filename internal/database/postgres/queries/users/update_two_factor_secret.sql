UPDATE users SET
	two_factor_secret_verified_at = $1,
	two_factor_secret = $2,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $3;
