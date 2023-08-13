// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: update_username.sql

package generated

import (
	"context"
)

const updateUserUsername = `-- name: UpdateUserUsername :exec

UPDATE users SET
	username = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2
`

type UpdateUserUsernameParams struct {
	Username string
	ID       string
}

func (q *Queries) UpdateUserUsername(ctx context.Context, db DBTX, arg *UpdateUserUsernameParams) error {
	_, err := db.ExecContext(ctx, updateUserUsername, arg.Username, arg.ID)
	return err
}
