-- name: CheckMealPlanGroceryListItemExistence :one

SELECT EXISTS ( SELECT meal_plan_grocery_list_items.id FROM meal_plan_grocery_list_items WHERE meal_plan_grocery_list_items.archived_at IS NULL AND meal_plan_grocery_list_items.id = $1 );
