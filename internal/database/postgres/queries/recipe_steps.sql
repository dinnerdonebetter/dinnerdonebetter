-- name: RecipeStepExists :one
SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_on IS NULL AND recipes.id = $3 );

-- name: GetRecipeStep :one
SELECT
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
    recipe_steps.explicit_instructions,
    recipe_steps.notes,
    recipe_steps.optional,
    recipe_steps.created_on,
    recipe_steps.last_updated_on,
    recipe_steps.archived_on,
    recipe_steps.belongs_to_recipe
FROM recipe_steps
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
         JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $1
  AND recipe_steps.id = $2
  AND recipes.archived_on IS NULL
  AND recipes.id = $3;

-- name: GetTotalRecipeStepsCount :one
SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL;

-- name: CreateRecipeStep :exec
INSERT INTO recipe_steps (id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,optional,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);

-- name: UpdateRecipeStep :exec
UPDATE recipe_steps SET
    index = $1,
    preparation_id = $2,
    minimum_estimated_time_in_seconds = $3,
    maximum_estimated_time_in_seconds = $4,
    minimum_temperature_in_celsius = $5,
    maximum_temperature_in_celsius = $6,
    explicit_instructions = $7,
    notes = $8,
    optional = $9,
    last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
  AND belongs_to_recipe = $10
  AND id = $11;

-- name: ArchiveRecipeStep :exec
UPDATE recipe_steps SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2;







