-- name: RecipeStepInstrumentExists :one
SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 );

-- name: GetRecipeStepInstrument :many
SELECT
    recipe_step_instruments.id,
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on,
    recipe_step_instruments.recipe_step_product_id,
    recipe_step_instruments.name,
    recipe_step_instruments.product_of_recipe_step,
    recipe_step_instruments.notes,
    recipe_step_instruments.preference_rank,
    recipe_step_instruments.created_on,
    recipe_step_instruments.last_updated_on,
    recipe_step_instruments.archived_on,
    recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
         LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
         JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_on IS NULL
  AND recipe_step_instruments.belongs_to_recipe_step = $1
  AND recipe_step_instruments.id = $2
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $3
  AND recipe_steps.id = $4
  AND recipes.archived_on IS NULL
  AND recipes.id = $5;

-- name: GetTotalRecipeStepInstrumentCount :one
SELECT COUNT(recipe_step_instruments.id) FROM recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL;

-- name: GetRecipeStepInstrumentsForRecipe :many
SELECT
    recipe_step_instruments.id,
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on,
    recipe_step_instruments.recipe_step_product_id,
    recipe_step_instruments.name,
    recipe_step_instruments.product_of_recipe_step,
    recipe_step_instruments.notes,
    recipe_step_instruments.preference_rank,
    recipe_step_instruments.created_on,
    recipe_step_instruments.last_updated_on,
    recipe_step_instruments.archived_on,
    recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
         LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
         JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_on IS NULL
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $1
  AND recipes.archived_on IS NULL
  AND recipes.id = $2;

-- name: CreateRecipeStepInstrument :exec
INSERT INTO recipe_step_instruments (id,instrument_id,recipe_step_product_id,name,product_of_recipe_step,notes,preference_rank,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);

-- name: UpdateRecipeStepInstrument :exec
UPDATE recipe_step_instruments
SET
    instrument_id = $1,
    recipe_step_product_id = $2,
    name = $3,
    product_of_recipe_step = $4,
    notes = $5,
    preference_rank = $6,
    last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
  AND belongs_to_recipe_step = $7
  AND id = $8;

-- name: ArchiveRecipeStepInstrument :exec
UPDATE recipe_step_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2;
