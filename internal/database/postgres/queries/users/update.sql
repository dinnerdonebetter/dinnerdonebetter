UPDATE users SET
	username = $1,
	hashed_password = $2,
	avatar_src = $3,
	two_factor_secret = $4,
	two_factor_secret_verified_at = $5,
	birth_day = $6,
	birth_month = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $8;
