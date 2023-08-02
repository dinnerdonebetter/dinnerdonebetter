-- name: UpdateUserDetails :exec

UPDATE users SET
	first_name = $1,
	last_name = $2,
	birthday = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $4;
