-- name: RecipeStepIngredientExists :one
SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 );

-- name: GetRecipeStepIngredient :many
SELECT
    recipe_step_ingredients.id,
    recipe_step_ingredients.name,
    recipe_step_ingredients.ingredient_id,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on,
    recipe_step_ingredients.minimum_quantity_value,
    recipe_step_ingredients.maximum_quantity_value,
    recipe_step_ingredients.quantity_notes,
    recipe_step_ingredients.product_of_recipe_step,
    recipe_step_ingredients.recipe_step_product_id,
    recipe_step_ingredients.ingredient_notes,
    recipe_step_ingredients.created_on,
    recipe_step_ingredients.last_updated_on,
    recipe_step_ingredients.archived_on,
    recipe_step_ingredients.belongs_to_recipe_step
FROM recipe_step_ingredients
  JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id
  JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
  JOIN valid_ingredients ON recipe_step_ingredients.ingredient_id=valid_ingredients.id
  JOIN valid_measurement_units ON recipe_step_ingredients.measurement_unit=valid_measurement_units.id
WHERE recipe_step_ingredients.archived_on IS NULL
  AND recipe_step_ingredients.belongs_to_recipe_step = $1
  AND recipe_step_ingredients.id = $2
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $3
  AND recipe_steps.id = $4
  AND recipes.archived_on IS NULL
  AND recipes.id = $5;

-- name: TotalRecipeStepIngredientCount :one
SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL;

-- name: CreateRecipeStepIngredient :exec
INSERT INTO recipe_step_ingredients (
    id,
    name,
    ingredient_id,
    measurement_unit,
    minimum_quantity_value,
    maximum_quantity_value,
    quantity_notes,
    product_of_recipe_step,
    recipe_step_product_id,
    ingredient_notes,
    belongs_to_recipe_step
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);

-- name: UpdateRecipeStepIngredient :exec

UPDATE recipe_step_ingredients SET
    ingredient_id = $1,
    name = $2,
    measurement_unit = $3,
    minimum_quantity_value = $4,
    maximum_quantity_value = $5,
    quantity_notes = $6,
    product_of_recipe_step = $7,
    recipe_step_product_id = $8,
    ingredient_notes = $9,
    last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL AND belongs_to_recipe_step = $10
  AND id = $11;

-- name: ArchiveRecipeStepIngredient :exec
UPDATE recipe_step_ingredients SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2;

