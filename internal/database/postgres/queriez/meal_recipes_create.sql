-- name: CreateMealRecipe :exec
INSERT INTO meal_recipes (id,meal_id,recipe_id) VALUES ($1,$2,$3);
