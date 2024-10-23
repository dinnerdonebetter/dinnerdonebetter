// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: maintenance.sql

package generated

import (
	"context"
)

const deleteExpiredOAuth2ClientTokens = `-- name: DeleteExpiredOAuth2ClientTokens :execrows
DELETE FROM oauth2_client_tokens WHERE code_expires_at < (NOW() - interval '1 day') AND access_expires_at < (NOW() - interval '1 day') AND refresh_expires_at < (NOW() - interval '1 day')
`

func (q *Queries) DeleteExpiredOAuth2ClientTokens(ctx context.Context, db DBTX) (int64, error) {
	result, err := db.ExecContext(ctx, deleteExpiredOAuth2ClientTokens)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}