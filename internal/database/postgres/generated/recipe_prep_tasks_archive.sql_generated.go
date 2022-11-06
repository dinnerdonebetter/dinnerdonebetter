// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_prep_tasks_archive.sql

package generated

import (
	"context"
)

const ArchiveRecipePrepTask = `-- name: ArchiveRecipePrepTask :exec
UPDATE recipe_prep_tasks SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveRecipePrepTask(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, ArchiveRecipePrepTask, id)
	return err
}
