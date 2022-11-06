// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_options_finalize.sql

package generated

import (
	"context"
	"database/sql"
)

const FinalizeMealPlanOption = `-- name: FinalizeMealPlanOption :exec
UPDATE meal_plan_options SET chosen = (belongs_to_meal_plan_event = $1 AND id = $2), tiebroken = $3 WHERE archived_at IS NULL AND belongs_to_meal_plan_event = $1 AND id = $2
`

type FinalizeMealPlanOptionParams struct {
	BelongsToMealPlanEvent sql.NullString
	ID                     string
	Tiebroken              bool
}

func (q *Queries) FinalizeMealPlanOption(ctx context.Context, arg *FinalizeMealPlanOptionParams) error {
	_, err := q.db.ExecContext(ctx, FinalizeMealPlanOption, arg.BelongsToMealPlanEvent, arg.ID, arg.Tiebroken)
	return err
}
