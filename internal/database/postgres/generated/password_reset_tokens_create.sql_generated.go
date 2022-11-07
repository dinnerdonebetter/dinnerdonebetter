// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: password_reset_tokens_create.sql

package generated

import (
	"context"
)

const CreatePasswordResetToken = `-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,NOW() + (30 * interval '1 minutes'),$3)
`

type CreatePasswordResetTokenParams struct {
	ID            string `db:"id"`
	Token         string `db:"token"`
	BelongsToUser string `db:"belongs_to_user"`
}

func (q *Queries) CreatePasswordResetToken(ctx context.Context, db DBTX, arg *CreatePasswordResetTokenParams) error {
	_, err := db.ExecContext(ctx, CreatePasswordResetToken, arg.ID, arg.Token, arg.BelongsToUser)
	return err
}