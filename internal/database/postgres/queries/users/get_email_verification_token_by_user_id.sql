SELECT
	users.email_address_verification_token
FROM users
WHERE users.archived_at IS NULL
    AND users.email_address_verified_at IS NULL
	AND users.id = $1;
