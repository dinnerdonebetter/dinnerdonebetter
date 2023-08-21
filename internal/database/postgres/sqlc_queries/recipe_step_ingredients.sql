-- name: ArchiveRecipeStepIngredient :exec

UPDATE recipe_step_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;


-- name: CreateRecipeStepIngredient :exec

INSERT INTO recipe_step_ingredients (
	id,
	"name",
	optional,
	ingredient_id,
	measurement_unit,
	minimum_quantity_value,
	maximum_quantity_value,
	quantity_notes,
	recipe_step_product_id,
	ingredient_notes,
	option_index,
	to_taste,
	product_percentage_to_use,
    vessel_index,
    recipe_step_product_recipe_id,
	belongs_to_recipe_step
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);


-- name: CheckRecipeStepIngredientExistence :one

SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_at IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = sqlc.arg(recipe_step_id) AND recipe_step_ingredients.id = sqlc.arg(recipe_step_ingredient_id) AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id) AND recipe_steps.id = sqlc.arg(recipe_step_id) AND recipes.archived_at IS NULL AND recipes.id = sqlc.arg(recipe_id) );


-- name: GetRecipeStepIngredientsForRecipe :many

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
    valid_ingredients.is_starch,
    valid_ingredients.is_protein,
    valid_ingredients.is_grain,
    valid_ingredients.is_fruit,
    valid_ingredients.is_salt,
    valid_ingredients.is_fat,
    valid_ingredients.is_acid,
    valid_ingredients.is_heat,
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
	recipe_step_ingredients.recipe_step_product_id,
	recipe_step_ingredients.ingredient_notes,
	recipe_step_ingredients.option_index,
	recipe_step_ingredients.to_taste,
	recipe_step_ingredients.product_percentage_to_use,
    recipe_step_ingredients.vessel_index,
    recipe_step_ingredients.recipe_step_product_recipe_id,
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

-- name: GetRecipeStepIngredient :one

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
    valid_ingredients.is_starch,
    valid_ingredients.is_protein,
    valid_ingredients.is_grain,
    valid_ingredients.is_fruit,
    valid_ingredients.is_salt,
    valid_ingredients.is_fat,
    valid_ingredients.is_acid,
    valid_ingredients.is_heat,
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
	recipe_step_ingredients.recipe_step_product_id,
	recipe_step_ingredients.ingredient_notes,
	recipe_step_ingredients.option_index,
	recipe_step_ingredients.to_taste,
	recipe_step_ingredients.product_percentage_to_use,
    recipe_step_ingredients.vessel_index,
    recipe_step_ingredients.recipe_step_product_recipe_id,
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


-- name: UpdateRecipeStepIngredient :exec

UPDATE recipe_step_ingredients SET
	ingredient_id = $1,
	name = $2,
	optional = $3,
	measurement_unit = $4,
	minimum_quantity_value = $5,
	maximum_quantity_value = $6,
	quantity_notes = $7,
	recipe_step_product_id = $8,
	ingredient_notes = $9,
	option_index = $10,
	to_taste = $11,
	product_percentage_to_use = $12,
    vessel_index = $13,
    recipe_step_product_recipe_id = $14,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND belongs_to_recipe_step = $15
	AND id = $16;
