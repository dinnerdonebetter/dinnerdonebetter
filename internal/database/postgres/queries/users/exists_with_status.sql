-- name: UserExistsWithStatus :one

SELECT EXISTS ( SELECT users.id FROM users WHERE users.archived_at IS NULL AND users.id = $1 AND users.user_account_status = $2);
