UPDATE users SET
	username = $1,
	first_name = $2,
	last_name = $3,
	hashed_password = $4,
	avatar_src = $5,
	birthday = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $8;
