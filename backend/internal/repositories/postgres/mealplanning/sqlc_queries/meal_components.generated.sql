-- name: CreateMealComponent :exec
INSERT INTO meal_components (
	id,
	belongs_to_meal,
	recipe_id,
	meal_component_type,
	recipe_scale
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_meal),
	sqlc.arg(recipe_id),
	sqlc.arg(meal_component_type),
	sqlc.arg(recipe_scale)
);
