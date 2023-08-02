// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: mark_two_factor_secret_as_verified.sql

package generated

import (
	"context"
)

const MarkTwoFactorSecretAsVerified = `-- name: MarkTwoFactorSecretAsVerified :exec

UPDATE users SET
	two_factor_secret_verified_at = NOW(),
	user_account_status = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $2
`

type MarkTwoFactorSecretAsVerifiedParams struct {
	UserAccountStatus string `db:"user_account_status"`
	ID                string `db:"id"`
}

func (q *Queries) MarkTwoFactorSecretAsVerified(ctx context.Context, db DBTX, arg *MarkTwoFactorSecretAsVerifiedParams) error {
	_, err := db.ExecContext(ctx, MarkTwoFactorSecretAsVerified, arg.UserAccountStatus, arg.ID)
	return err
}