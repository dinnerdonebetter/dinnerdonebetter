-- name: RecipeStepProductExists :one
SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 );

-- name: GetRecipeStepProduct :many
SELECT
    recipe_step_products.id,
    recipe_step_products.name,
    recipe_step_products.type,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on,
    recipe_step_products.minimum_quantity_value,
    recipe_step_products.maximum_quantity_value,
    recipe_step_products.quantity_notes,
    recipe_step_products.created_on,
    recipe_step_products.last_updated_on,
    recipe_step_products.archived_on,
    recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
         JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
         JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_on IS NULL
  AND recipe_step_products.belongs_to_recipe_step = $1
  AND recipe_step_products.id = $2
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $3
  AND recipe_steps.id = $4
  AND recipes.archived_on IS NULL
  AND recipes.id = $5;

-- name: TotalRecipeStepProductCount :one
SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL;

-- name: RecipeStepProductsForRecipeQuery :many
SELECT
    recipe_step_products.id,
    recipe_step_products.name,
    recipe_step_products.type,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on,
    recipe_step_products.minimum_quantity_value,
    recipe_step_products.maximum_quantity_value,
    recipe_step_products.quantity_notes,
    recipe_step_products.created_on,
    recipe_step_products.last_updated_on,
    recipe_step_products.archived_on,
    recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
         JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
         LEFT OUTER JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_on IS NULL
  AND recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $1
  AND recipes.archived_on IS NULL
  AND recipes.id = $2;

-- name: CreateRecipeStepProduct :exec
INSERT INTO recipe_step_products (id,name,type,measurement_unit,minimum_quantity_value,maximum_quantity_value,quantity_notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);

-- name: UpdateRecipeStepProduct :exec
UPDATE recipe_step_products SET name = $1, type = $2, measurement_unit = $3, minimum_quantity_value = $4, maximum_quantity_value = $5, quantity_notes = $6, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $7 AND id = $8;

-- name: ArchiveRecipeStepProduct :exec
UPDATE recipe_step_products SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2;
