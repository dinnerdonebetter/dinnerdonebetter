-- name: ArchiveMealPlanGroceryListItem :exec

UPDATE meal_plan_grocery_list_items SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
