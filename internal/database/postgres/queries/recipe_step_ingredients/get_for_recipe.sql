SELECT
	recipe_step_ingredients.id,
	recipe_step_ingredients.name,
	recipe_step_ingredients.optional,
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.volumetric,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
    valid_ingredients.storage_instructions,
    valid_ingredients.slug,
    valid_ingredients.contains_alcohol,
    valid_ingredients.shopping_suggestions,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
    valid_measurement_units.imperial,
    valid_measurement_units.slug,
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
    recipe_step_ingredients.option_index,
    recipe_step_ingredients.requires_defrost,
	recipe_step_ingredients.created_at,
	recipe_step_ingredients.last_updated_at,
	recipe_step_ingredients.archived_at,
	recipe_step_ingredients.belongs_to_recipe_step
FROM
	recipe_step_ingredients
	JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	LEFT JOIN valid_ingredients ON recipe_step_ingredients.ingredient_id = valid_ingredients.id
	JOIN valid_measurement_units ON recipe_step_ingredients.measurement_unit = valid_measurement_units.id
WHERE
	recipe_step_ingredients.archived_at IS NULL
	AND recipes.id = $1
GROUP BY
	recipe_step_ingredients.id,
	valid_measurement_units.id,
	valid_ingredients.id
ORDER BY
	recipe_step_ingredients.id;