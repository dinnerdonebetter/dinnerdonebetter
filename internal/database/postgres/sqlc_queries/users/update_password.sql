UPDATE users SET
	hashed_password = $1,
	requires_password_change = $2,
	password_last_changed_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $3;
