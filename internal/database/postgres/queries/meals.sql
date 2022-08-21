-- name: MealExists :one
SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_on IS NULL AND meals.id = $1 );

-- name: GetMealByID :one
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.created_on,
	meals.last_updated_on,
	meals.archived_on,
	meals.created_by_user,
	meal_recipes.recipe_id
FROM meals
	FULL OUTER JOIN meal_recipes ON meal_recipes.meal_id=meals.id
WHERE meals.archived_on IS NULL
	AND meal_recipes.archived_on IS NULL
	AND meals.id = $1;

-- name: GetTotalMealsCount :one
SELECT COUNT(meals.id) FROM meals WHERE meals.archived_on IS NULL;

-- name: CreateMeal :exec
INSERT INTO meals (id,name,description,created_by_user) VALUES ($1,$2,$3,$4);

-- name: CreateMealRecipe :exec
INSERT INTO meal_recipes (id,meal_id,recipe_id) VALUES ($1,$2,$3);

-- name: ArchiveMeal :exec
UPDATE meals SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND created_by_user = $1 AND id = $2;
