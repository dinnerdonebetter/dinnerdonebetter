// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: add_user_to_household.sql

package generated

import (
	"context"
)

const addUserToHousehold = `-- name: AddUserToHousehold :exec

INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_role)
VALUES ($1,$2,$3,$4)
`

type AddUserToHouseholdParams struct {
	ID                 string
	BelongsToUser      string
	BelongsToHousehold string
	HouseholdRole      string
}

func (q *Queries) AddUserToHousehold(ctx context.Context, db DBTX, arg *AddUserToHouseholdParams) error {
	_, err := db.ExecContext(ctx, addUserToHousehold,
		arg.ID,
		arg.BelongsToUser,
		arg.BelongsToHousehold,
		arg.HouseholdRole,
	)
	return err
}
