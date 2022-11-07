// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plans_create.sql

package generated

import (
	"context"
	"time"
)

const CreateMealPlan = `-- name: CreateMealPlan :exec
INSERT INTO meal_plans (id,notes,status,voting_deadline,belongs_to_household) VALUES ($1,$2,$3,$4,$5)
`

type CreateMealPlanParams struct {
	ID                 string         `db:"id"`
	Notes              string         `db:"notes"`
	Status             MealPlanStatus `db:"status"`
	VotingDeadline     time.Time      `db:"voting_deadline"`
	BelongsToHousehold string         `db:"belongs_to_household"`
}

func (q *Queries) CreateMealPlan(ctx context.Context, db DBTX, arg *CreateMealPlanParams) error {
	_, err := db.ExecContext(ctx, CreateMealPlan,
		arg.ID,
		arg.Notes,
		arg.Status,
		arg.VotingDeadline,
		arg.BelongsToHousehold,
	)
	return err
}