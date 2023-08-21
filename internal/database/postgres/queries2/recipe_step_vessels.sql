-- name: ArchiveRecipeStepVessel :exec

UPDATE recipe_step_vessels SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;


-- name: CreateRecipeStepVessel :exec

INSERT INTO recipe_step_vessels
(id,"name",notes,belongs_to_recipe_step,recipe_step_product_id,valid_vessel_id,vessel_predicate,minimum_quantity,maximum_quantity,unavailable_after_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);


-- name: CheckRecipeStepVesselExistence :one

SELECT EXISTS (
    SELECT recipe_step_vessels.id
    FROM recipe_step_vessels
        JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
        JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    WHERE recipe_step_vessels.archived_at IS NULL
        AND recipe_step_vessels.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
        AND recipe_step_vessels.id = sqlc.arg(recipe_step_vessel_id)
        AND recipe_steps.archived_at IS NULL
        AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
        AND recipe_steps.id = sqlc.arg(recipe_step_id)
        AND recipes.archived_at IS NULL
        AND recipes.id = sqlc.arg(recipe_id)
);


-- name: GetRecipeStepVesselsForRecipe :many

SELECT
	recipe_step_vessels.id,
	valid_vessels.id,
    valid_vessels.name,
    valid_vessels.plural_name,
    valid_vessels.description,
    valid_vessels.icon_path,
    valid_vessels.usable_for_storage,
    valid_vessels.slug,
    valid_vessels.display_in_summary_lists,
    valid_vessels.include_in_generated_instructions,
    valid_vessels.capacity,
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
    valid_vessels.width_in_millimeters,
    valid_vessels.length_in_millimeters,
    valid_vessels.height_in_millimeters,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at,
    recipe_step_vessels.name,
    recipe_step_vessels.notes,
    recipe_step_vessels.belongs_to_recipe_step,
    recipe_step_vessels.recipe_step_product_id,
    recipe_step_vessels.vessel_predicate,
    recipe_step_vessels.minimum_quantity,
    recipe_step_vessels.maximum_quantity,
    recipe_step_vessels.unavailable_after_step,
    recipe_step_vessels.created_at,
    recipe_step_vessels.last_updated_at,
    recipe_step_vessels.archived_at
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
    LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;


-- name: GetRecipeStepVessel :one

SELECT
    recipe_step_vessels.id,
	valid_vessels.id,
    valid_vessels.name,
    valid_vessels.plural_name,
    valid_vessels.description,
    valid_vessels.icon_path,
    valid_vessels.usable_for_storage,
    valid_vessels.slug,
    valid_vessels.display_in_summary_lists,
    valid_vessels.include_in_generated_instructions,
    valid_vessels.capacity,
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
    valid_vessels.width_in_millimeters,
    valid_vessels.length_in_millimeters,
    valid_vessels.height_in_millimeters,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at,
    recipe_step_vessels.name,
    recipe_step_vessels.notes,
    recipe_step_vessels.belongs_to_recipe_step,
    recipe_step_vessels.recipe_step_product_id,
    recipe_step_vessels.vessel_predicate,
    recipe_step_vessels.minimum_quantity,
    recipe_step_vessels.maximum_quantity,
    recipe_step_vessels.unavailable_after_step,
    recipe_step_vessels.created_at,
    recipe_step_vessels.last_updated_at,
    recipe_step_vessels.archived_at
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
    LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_step_vessels.belongs_to_recipe_step = $1
	AND recipe_step_vessels.id = $2
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $3
	AND recipe_steps.id = $4
	AND recipes.archived_at IS NULL
	AND recipes.id = $5;


-- name: UpdateRecipeStepVessel :exec

UPDATE recipe_step_vessels SET
	name = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	recipe_step_product_id = $4,
	valid_vessel_id = $5,
	vessel_predicate = $6,
	minimum_quantity = $7,
    maximum_quantity = $8,
    unavailable_after_step = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $10
	AND id = $11;