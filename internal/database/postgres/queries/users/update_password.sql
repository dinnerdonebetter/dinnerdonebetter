-- name: UpdateUserPassword :exec

UPDATE users SET
	hashed_password = $1,
	password_last_changed_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
