UPDATE users SET
	two_factor_secret_verified_at = NOW(),
	user_account_status = $1
WHERE archived_at IS NULL
	AND id = $2;
