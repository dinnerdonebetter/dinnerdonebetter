-- name: CreatePasswordResetToken :exec

INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,NOW() + (30 * interval '1 minutes'),$3);


-- name: GetPasswordResetToken :one

SELECT
	password_reset_tokens.id,
	password_reset_tokens.token,
	password_reset_tokens.expires_at,
	password_reset_tokens.created_at,
	password_reset_tokens.last_updated_at,
	password_reset_tokens.redeemed_at,
	password_reset_tokens.belongs_to_user
FROM password_reset_tokens
WHERE password_reset_tokens.redeemed_at IS NULL
	AND NOW() < password_reset_tokens.expires_at
	AND password_reset_tokens.token = $1;


-- name: RedeemPasswordResetToken :exec

UPDATE password_reset_tokens SET redeemed_at = NOW() WHERE redeemed_at IS NULL AND id = $1;
