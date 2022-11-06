// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: users_archive.sql

package generated

import (
	"context"
)

const ArchiveUser = `-- name: ArchiveUser :exec
UPDATE users SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) ArchiveUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, ArchiveUser, id)
	return err
}
