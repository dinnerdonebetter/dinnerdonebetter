-- name: ArchiveValidIngredientMeasurementUnit :exec

UPDATE valid_ingredient_measurement_units SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateValidIngredientMeasurementUnit :exec

INSERT INTO valid_ingredient_measurement_units
(id,notes,valid_measurement_unit_id,valid_ingredient_id,minimum_allowable_quantity,maximum_allowable_quantity)
VALUES ($1,$2,$3,$4,$5,$6);

-- name: CheckValidIngredientMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_at IS NULL AND valid_ingredient_measurement_units.id = $1 );

-- name: GetValidIngredientMeasurementUnitsForIngredient :many

SELECT
	valid_ingredient_measurement_units.id as valid_ingredient_measurement_unit_id,
	valid_ingredient_measurement_units.notes as valid_ingredient_measurement_unit_notes,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	valid_ingredient_measurement_units.minimum_allowable_quantity as valid_ingredient_measurement_unit_minimum_allowable_quantity,
	valid_ingredient_measurement_units.maximum_allowable_quantity as valid_ingredient_measurement_unit_maximum_allowable_quantity,
	valid_ingredient_measurement_units.created_at as valid_ingredient_measurement_unit_created_at,
	valid_ingredient_measurement_units.last_updated_at as valid_ingredient_measurement_unit_last_updated_at,
	valid_ingredient_measurement_units.archived_at as valid_ingredient_measurement_unit_archived_at,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
          AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
    ) as total_count
FROM valid_ingredient_measurement_units
	JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    AND valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetValidIngredientMeasurementUnitsForMeasurementUnit :many

SELECT
	valid_ingredient_measurement_units.id as valid_ingredient_measurement_unit_id,
	valid_ingredient_measurement_units.notes as valid_ingredient_measurement_unit_notes,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	valid_ingredient_measurement_units.minimum_allowable_quantity as valid_ingredient_measurement_unit_minimum_allowable_quantity,
	valid_ingredient_measurement_units.maximum_allowable_quantity as valid_ingredient_measurement_unit_maximum_allowable_quantity,
	valid_ingredient_measurement_units.created_at as valid_ingredient_measurement_unit_created_at,
	valid_ingredient_measurement_units.last_updated_at as valid_ingredient_measurement_unit_last_updated_at,
	valid_ingredient_measurement_units.archived_at as valid_ingredient_measurement_unit_archived_at,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
          AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
    ) as total_count
FROM valid_ingredient_measurement_units
	JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    AND valid_ingredient_measurement_units.valid_measurement_unit_id = sqlc.arg(valid_measurement_unit_id)
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetValidIngredientMeasurementUnits :many

SELECT
	valid_ingredient_measurement_units.id as valid_ingredient_measurement_unit_id,
	valid_ingredient_measurement_units.notes as valid_ingredient_measurement_unit_notes,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	valid_ingredient_measurement_units.minimum_allowable_quantity as valid_ingredient_measurement_unit_minimum_allowable_quantity,
	valid_ingredient_measurement_units.maximum_allowable_quantity as valid_ingredient_measurement_unit_maximum_allowable_quantity,
	valid_ingredient_measurement_units.created_at as valid_ingredient_measurement_unit_created_at,
	valid_ingredient_measurement_units.last_updated_at as valid_ingredient_measurement_unit_last_updated_at,
	valid_ingredient_measurement_units.archived_at as valid_ingredient_measurement_unit_archived_at,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
          AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_measurement_units.id)
        FROM
            valid_ingredient_measurement_units
        WHERE
            valid_ingredient_measurement_units.archived_at IS NULL
    ) as total_count
FROM valid_ingredient_measurement_units
	JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND valid_ingredient_measurement_units.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND valid_ingredient_measurement_units.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (valid_ingredient_measurement_units.last_updated_at IS NULL OR valid_ingredient_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetValidIngredientMeasurementUnit :one

SELECT
	valid_ingredient_measurement_units.id as valid_ingredient_measurement_unit_id,
	valid_ingredient_measurement_units.notes as valid_ingredient_measurement_unit_notes,
	valid_measurement_units.id as valid_measurement_unit_id,
	valid_measurement_units.name as valid_measurement_unit_name,
	valid_measurement_units.description as valid_measurement_unit_description,
	valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
	valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
	valid_measurement_units.universal as valid_measurement_unit_universal,
	valid_measurement_units.metric as valid_measurement_unit_metric,
	valid_measurement_units.imperial as valid_measurement_unit_imperial,
	valid_measurement_units.slug as valid_measurement_unit_slug,
	valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
	valid_measurement_units.created_at as valid_measurement_unit_created_at,
	valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
	valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
    valid_ingredients.is_starch as valid_ingredient_is_starch,
    valid_ingredients.is_protein as valid_ingredient_is_protein,
    valid_ingredients.is_grain as valid_ingredient_is_grain,
    valid_ingredients.is_fruit as valid_ingredient_is_fruit,
    valid_ingredients.is_salt as valid_ingredient_is_salt,
    valid_ingredients.is_fat as valid_ingredient_is_fat,
    valid_ingredients.is_acid as valid_ingredient_is_acid,
    valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	valid_ingredient_measurement_units.minimum_allowable_quantity as valid_ingredient_measurement_unit_minimum_allowable_quantity,
	valid_ingredient_measurement_units.maximum_allowable_quantity as valid_ingredient_measurement_unit_maximum_allowable_quantity,
	valid_ingredient_measurement_units.created_at as valid_ingredient_measurement_unit_created_at,
	valid_ingredient_measurement_units.last_updated_at as valid_ingredient_measurement_unit_last_updated_at,
	valid_ingredient_measurement_units.archived_at as valid_ingredient_measurement_unit_archived_at
FROM valid_ingredient_measurement_units
	JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_at IS NULL
	AND valid_ingredient_measurement_units.id = $1;

-- name: ValidIngredientMeasurementUnitPairIsValid :one

SELECT EXISTS(
	SELECT id
	FROM valid_ingredient_measurement_units
	WHERE valid_measurement_unit_id = $1
	AND valid_ingredient_id = $2
	AND archived_at IS NULL
);

-- name: UpdateValidIngredientMeasurementUnit :exec

UPDATE valid_ingredient_measurement_units
SET
	notes = $1,
	valid_measurement_unit_id = $2,
	valid_ingredient_id = $3,
	minimum_allowable_quantity = $4,
	maximum_allowable_quantity = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
