// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plans_mark_as_grocery_list_initialized.sql

package generated

import (
	"context"
)

const MarkMealPlanAsHavingGroceryListInitialized = `-- name: MarkMealPlanAsHavingGroceryListInitialized :exec
UPDATE meal_plans
SET
    grocery_list_initialized = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, MarkMealPlanAsHavingGroceryListInitialized, id)
	return err
}
