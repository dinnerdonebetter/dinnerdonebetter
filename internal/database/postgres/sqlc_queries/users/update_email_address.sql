-- name: UpdateUserEmailAddress :exec

UPDATE users SET
	email_address = $1,
	email_address_verified_at = NULL,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2;
