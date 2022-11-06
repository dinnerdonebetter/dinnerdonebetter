// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: household_invitations_create.sql

package generated

import (
	"context"
	"database/sql"
)

const CreateHouseholdInvitation = `-- name: CreateHouseholdInvitation :exec
INSERT INTO household_invitations (id,from_user,to_user,note,to_email,token,destination_household) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

type CreateHouseholdInvitationParams struct {
	ID                   string
	FromUser             string
	ToUser               sql.NullString
	Note                 string
	ToEmail              string
	Token                string
	DestinationHousehold string
}

func (q *Queries) CreateHouseholdInvitation(ctx context.Context, arg *CreateHouseholdInvitationParams) error {
	_, err := q.db.ExecContext(ctx, CreateHouseholdInvitation,
		arg.ID,
		arg.FromUser,
		arg.ToUser,
		arg.Note,
		arg.ToEmail,
		arg.Token,
		arg.DestinationHousehold,
	)
	return err
}
