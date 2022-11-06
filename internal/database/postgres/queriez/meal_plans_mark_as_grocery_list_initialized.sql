-- name: MarkMealPlanAsHavingGroceryListInitialized :exec
UPDATE meal_plans
SET
    grocery_list_initialized = 'true',
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;
