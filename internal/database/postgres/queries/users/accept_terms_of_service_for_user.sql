UPDATE users SET
	last_accepted_terms_of_service = NOW()
WHERE archived_at IS NULL
	AND id = $1;
