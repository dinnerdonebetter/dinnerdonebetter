// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_events_create.sql

package generated

import (
	"context"
	"time"
)

const CreateMealPlanEvent = `-- name: CreateMealPlanEvent :exec
INSERT INTO
	meal_plan_events (id, notes, starts_at, ends_at, meal_name, belongs_to_meal_plan)
VALUES
	($1, $2, $3, $4, $5, $6)
`

type CreateMealPlanEventParams struct {
	ID                string
	Notes             string
	StartsAt          time.Time
	EndsAt            time.Time
	MealName          MealName
	BelongsToMealPlan string
}

func (q *Queries) CreateMealPlanEvent(ctx context.Context, arg *CreateMealPlanEventParams) error {
	_, err := q.db.ExecContext(ctx, CreateMealPlanEvent,
		arg.ID,
		arg.Notes,
		arg.StartsAt,
		arg.EndsAt,
		arg.MealName,
		arg.BelongsToMealPlan,
	)
	return err
}
