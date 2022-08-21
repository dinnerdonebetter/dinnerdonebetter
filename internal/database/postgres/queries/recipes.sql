-- name: RecipeExists :one
SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_on IS NULL AND recipes.id = $1 );

-- name: GetRecipeByID :many
SELECT
    recipes.id,
    recipes.name,
    recipes.source,
    recipes.description,
    recipes.inspired_by_recipe_id,
    recipes.created_on,
    recipes.last_updated_on,
    recipes.archived_on,
    recipes.created_by_user,
    recipe_steps.id,
    recipe_steps.index,
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on,
    recipe_steps.minimum_estimated_time_in_seconds,
    recipe_steps.maximum_estimated_time_in_seconds,
    recipe_steps.minimum_temperature_in_celsius,
    recipe_steps.maximum_temperature_in_celsius,
    recipe_steps.notes,
    recipe_steps.optional,
    recipe_steps.created_on,
    recipe_steps.last_updated_on,
    recipe_steps.archived_on,
    recipe_steps.belongs_to_recipe
FROM recipes
         FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
         FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_on IS NULL
  AND recipes.id = $1
ORDER BY recipe_steps.index;

-- name: GetRecipeByIDAndAuthorID :many
SELECT
    recipes.id,
    recipes.name,
    recipes.source,
    recipes.description,
    recipes.inspired_by_recipe_id,
    recipes.created_on,
    recipes.last_updated_on,
    recipes.archived_on,
    recipes.created_by_user,
    recipe_steps.id,
    recipe_steps.index,
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on,
    recipe_steps.minimum_estimated_time_in_seconds,
    recipe_steps.maximum_estimated_time_in_seconds,
    recipe_steps.minimum_temperature_in_celsius,
    recipe_steps.maximum_temperature_in_celsius,
    recipe_steps.notes,
    recipe_steps.optional,
    recipe_steps.created_on,
    recipe_steps.last_updated_on,
    recipe_steps.archived_on,
    recipe_steps.belongs_to_recipe
FROM recipes
         FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
         FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_on IS NULL
  AND recipes.id = $1
  AND recipes.created_by_user = $2
ORDER BY recipe_steps.index;

-- name: GetTotalRecipesCount :one
SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL;

-- name: CreateRecipe :exec
INSERT INTO recipes (id,name,source,description,inspired_by_recipe_id,created_by_user) VALUES ($1,$2,$3,$4,$5,$6);

-- name: UpdateRecipe :exec
UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND created_by_user = $5 AND id = $6;

-- name: ArchiveRecipe :exec
UPDATE recipes SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND created_by_user = $1 AND id = $2;
