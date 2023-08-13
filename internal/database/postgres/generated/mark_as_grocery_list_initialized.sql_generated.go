// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: mark_as_grocery_list_initialized.sql

package generated

import (
	"context"
)

const markMealPlanAsGroceryListInitialized = `-- name: MarkMealPlanAsGroceryListInitialized :exec

UPDATE meal_plans
SET
	grocery_list_initialized = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) MarkMealPlanAsGroceryListInitialized(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, markMealPlanAsGroceryListInitialized, id)
	return err
}
