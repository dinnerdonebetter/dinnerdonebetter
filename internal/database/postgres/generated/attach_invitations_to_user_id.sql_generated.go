// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: attach_invitations_to_user_id.sql

package generated

import (
	"context"
	"database/sql"
)

const AttachHouseholdInvitationsToUserID = `-- name: AttachHouseholdInvitationsToUserID :exec

UPDATE household_invitations SET
	to_user = $1,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND to_email = LOWER($2)
`

type AttachHouseholdInvitationsToUserIDParams struct {
	ToUser sql.NullString `db:"to_user"`
	Lower  string         `db:"lower"`
}

func (q *Queries) AttachHouseholdInvitationsToUserID(ctx context.Context, db DBTX, arg *AttachHouseholdInvitationsToUserIDParams) error {
	_, err := db.ExecContext(ctx, AttachHouseholdInvitationsToUserID, arg.ToUser, arg.Lower)
	return err
}
