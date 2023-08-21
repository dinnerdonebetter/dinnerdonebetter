-- name: MarkTwoFactorSecretAsVerified :exec

UPDATE users SET
	two_factor_secret_verified_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;
