-- name: SetUserAccountStatus :exec
UPDATE users SET user_account_status = $1, user_account_status_explanation = $2 WHERE archived_on IS NULL AND id = $3;
