-- name: GetMeal :exec
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_recipes.recipe_id
FROM meals
	FULL OUTER JOIN meal_recipes ON meal_recipes.meal_id=meals.id
WHERE meals.archived_at IS NULL
	AND meal_recipes.archived_at IS NULL
	AND meals.id = $1;
