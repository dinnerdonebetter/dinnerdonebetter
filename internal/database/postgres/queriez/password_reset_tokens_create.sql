-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,NOW() + (30 * interval '1 minutes'),$3);
