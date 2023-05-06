UPDATE users SET
	username = $1,
	first_name = $2,
	last_name = $3,
	hashed_password = $4,
	avatar_src = $5,
	two_factor_secret = $6,
	two_factor_secret_verified_at = $7,
	birthday = $8,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $9;
