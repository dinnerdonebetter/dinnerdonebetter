// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: create_for_new_user.sql

package generated

import (
	"context"
)

const CreateHouseholdUserMembershipForNewUser = `-- name: CreateHouseholdUserMembershipForNewUser :exec

INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,default_household,household_role)
VALUES ($1,$2,$3,$4,$5)
`

type CreateHouseholdUserMembershipForNewUserParams struct {
	ID                 string `db:"id"`
	BelongsToUser      string `db:"belongs_to_user"`
	BelongsToHousehold string `db:"belongs_to_household"`
	DefaultHousehold   bool   `db:"default_household"`
	HouseholdRole      string `db:"household_role"`
}

func (q *Queries) CreateHouseholdUserMembershipForNewUser(ctx context.Context, db DBTX, arg *CreateHouseholdUserMembershipForNewUserParams) error {
	_, err := db.ExecContext(ctx, CreateHouseholdUserMembershipForNewUser,
		arg.ID,
		arg.BelongsToUser,
		arg.BelongsToHousehold,
		arg.DefaultHousehold,
		arg.HouseholdRole,
	)
	return err
}
