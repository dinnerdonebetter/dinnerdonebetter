-- name: CheckMealPlanGroceryListItemExistence :one

SELECT EXISTS ( SELECT meal_plan_grocery_list_items.id FROM meal_plan_grocery_list_items WHERE meal_plan_grocery_list_items.archived_at IS NULL AND meal_plan_grocery_list_items.id = sqlc.arg(meal_plan_grocery_list_item_id) AND meal_plan_grocery_list_items.belongs_to_meal_plan = sqlc.arg(meal_plan_id) );
