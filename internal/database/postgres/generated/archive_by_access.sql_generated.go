// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: archive_by_access.sql

package generated

import (
	"context"
)

const ArchiveOAuth2ClientTokenByAccess = `-- name: ArchiveOAuth2ClientTokenByAccess :exec

DELETE FROM oauth2_client_tokens WHERE access = $1
`

func (q *Queries) ArchiveOAuth2ClientTokenByAccess(ctx context.Context, db DBTX, access string) error {
	_, err := db.ExecContext(ctx, ArchiveOAuth2ClientTokenByAccess, access)
	return err
}
