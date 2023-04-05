SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
    meals.max_estimated_portions,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.recipe_id,
	meal_components.recipe_scale,
	meal_components.meal_component_type
FROM meals
	FULL OUTER JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
	AND meal_components.archived_at IS NULL
	AND meals.id = $1;
