-- name: SetUserAccountStatus :execresult

UPDATE users SET user_account_status = $1, user_account_status_explanation = $2 WHERE archived_at IS NULL AND id = $3;