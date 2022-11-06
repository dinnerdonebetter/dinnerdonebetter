-- name: ArchiveMealPlanGroceryListItem :exec
UPDATE meal_plan_grocery_list_items SET completed_at = NOW() WHERE completed_at IS NULL AND id = $1;
