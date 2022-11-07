// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plans_finalize.sql

package generated

import (
	"context"
)

const FinalizeMealPlan = `-- name: FinalizeMealPlan :exec
UPDATE meal_plans SET status = $1 WHERE archived_at IS NULL AND id = $2
`

type FinalizeMealPlanParams struct {
	Status MealPlanStatus `db:"status"`
	ID     string         `db:"id"`
}

func (q *Queries) FinalizeMealPlan(ctx context.Context, db DBTX, arg *FinalizeMealPlanParams) error {
	_, err := db.ExecContext(ctx, FinalizeMealPlan, arg.Status, arg.ID)
	return err
}