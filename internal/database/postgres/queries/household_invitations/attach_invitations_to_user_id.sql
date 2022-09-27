UPDATE household_invitations SET
	to_user = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND to_email = LOWER($2);
