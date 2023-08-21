-- name: ArchiveRecipeStepProduct :exec

UPDATE recipe_step_products SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;


-- name: CreateRecipeStepProduct :exec

INSERT INTO recipe_step_products
(id,"name","type",measurement_unit,minimum_quantity_value,maximum_quantity_value,quantity_notes,compostable,maximum_storage_duration_in_seconds,minimum_storage_temperature_in_celsius,maximum_storage_temperature_in_celsius,storage_instructions,belongs_to_recipe_step,is_liquid,is_waste,"index",contained_in_vessel_index)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17);


-- name: CheckRecipeStepProductExistence :one

SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_at IS NULL AND recipe_step_products.belongs_to_recipe_step = sqlc.arg(recipe_step_id) AND recipe_step_products.id = sqlc.arg(recipe_step_product_id) AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id) AND recipe_steps.id = sqlc.arg(recipe_step_id) AND recipes.archived_at IS NULL AND recipes.id = sqlc.arg(recipe_id) );


-- name: GetRecipeStepProductsForRecipe :many

SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.is_liquid,
	recipe_step_products.is_waste,
    recipe_step_products.index,
    recipe_step_products.contained_in_vessel_index,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	LEFT OUTER JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;


-- name: GetRecipeStepProduct :one

SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.is_liquid,
	recipe_step_products.is_waste,
    recipe_step_products.index,
    recipe_step_products.contained_in_vessel_index,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_step_products.belongs_to_recipe_step = $1
	AND recipe_step_products.id = $2
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $3
	AND recipe_steps.id = $4
	AND recipes.archived_at IS NULL
	AND recipes.id = $5;


-- name: UpdateRecipeStepProduct :exec

UPDATE recipe_step_products
SET
	"name" = $1,
	"type" = $2,
	measurement_unit = $3,
	minimum_quantity_value = $4,
	maximum_quantity_value = $5,
	quantity_notes = $6,
	compostable = $7,
	maximum_storage_duration_in_seconds = $8,
	minimum_storage_temperature_in_celsius = $9,
	maximum_storage_temperature_in_celsius = $10,
	storage_instructions = $11,
	is_liquid = $12,
	is_waste = $13,
    "index" = $14,
    contained_in_vessel_index = $15,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $16
	AND id = $17;
