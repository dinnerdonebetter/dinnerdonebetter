UPDATE users SET
	avatar_src = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
