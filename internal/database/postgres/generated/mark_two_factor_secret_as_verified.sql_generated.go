// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: mark_two_factor_secret_as_verified.sql

package generated

import (
	"context"
)

const markTwoFactorSecretAsVerified = `-- name: MarkTwoFactorSecretAsVerified :exec

UPDATE users SET
	two_factor_secret_verified_at = NOW(),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) MarkTwoFactorSecretAsVerified(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, markTwoFactorSecretAsVerified, id)
	return err
}
