UPDATE users SET
	email_address_verified_at = NOW()
WHERE email_address_verified_at IS NULL
	AND id = $1
	AND email_address_verification_token = $2;
