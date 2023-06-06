UPDATE users SET
	last_accepted_privacy_policy = NOW()
WHERE archived_at IS NULL
	AND id = $1;
