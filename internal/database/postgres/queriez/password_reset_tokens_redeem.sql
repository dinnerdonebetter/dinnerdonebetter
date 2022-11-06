-- name: RedeemPasswordResetToken :exec
UPDATE password_reset_tokens SET redeemed_at = NOW() WHERE redeemed_at IS NULL AND id = $1;
