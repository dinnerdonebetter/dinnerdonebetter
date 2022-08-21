-- name: GetPasswordResetToken :one
SELECT
	password_reset_tokens.id,
	password_reset_tokens.token,
	password_reset_tokens.expires_at,
	password_reset_tokens.created_on,
	password_reset_tokens.last_updated_on,
	password_reset_tokens.redeemed_on,
	password_reset_tokens.belongs_to_user
FROM password_reset_tokens
WHERE password_reset_tokens.redeemed_on IS NULL
AND extract(epoch from NOW()) < password_reset_tokens.expires_at
AND password_reset_tokens.token = $1;

-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,extract(epoch from (NOW() + (30 * interval '1 minutes'))),$3);

-- name: RedeemPasswordResetToken :exec
UPDATE password_reset_tokens SET redeemed_on = extract(epoch FROM NOW()) WHERE redeemed_on IS NULL AND id = $1;
