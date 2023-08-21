-- name: CreateMealComponent :exec

INSERT INTO meal_components (id,meal_id,recipe_id,meal_component_type,recipe_scale) VALUES ($1,$2,$3,$4,$5);
