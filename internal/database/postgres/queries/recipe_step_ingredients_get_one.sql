SELECT
    recipe_step_ingredients.id,
    recipe_step_ingredients.name,
    recipe_step_ingredients.optional,
    recipe_step_ingredients.ingredient_id,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.universal,
    valid_measurement_units.metric,
    valid_measurement_units.imperial,
    valid_measurement_units.plural_name,
    valid_measurement_units.created_at,
    valid_measurement_units.last_updated_at,
    valid_measurement_units.archived_at,
    recipe_step_ingredients.minimum_quantity_value,
    recipe_step_ingredients.maximum_quantity_value,
    recipe_step_ingredients.quantity_notes,
    recipe_step_ingredients.product_of_recipe_step,
    recipe_step_ingredients.recipe_step_product_id,
    recipe_step_ingredients.ingredient_notes,
    recipe_step_ingredients.created_at,
    recipe_step_ingredients.last_updated_at,
    recipe_step_ingredients.archived_at,
    recipe_step_ingredients.belongs_to_recipe_step
FROM recipe_step_ingredients
         JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
         JOIN valid_ingredients ON recipe_step_ingredients.ingredient_id=valid_ingredients.id
         JOIN valid_measurement_units ON recipe_step_ingredients.measurement_unit=valid_measurement_units.id
WHERE recipe_step_ingredients.archived_at IS NULL
  AND recipe_step_ingredients.belongs_to_recipe_step = $1
  AND recipe_step_ingredients.id = $2
  AND recipe_steps.archived_at IS NULL
  AND recipe_steps.belongs_to_recipe = $3
  AND recipe_steps.id = $4
  AND recipes.archived_at IS NULL
  AND recipes.id = $5;
