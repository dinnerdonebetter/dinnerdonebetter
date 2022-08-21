-- name: ValidIngredientMeasurementUnitExists :one
SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL AND valid_ingredient_measurement_units.id = $1 );

-- name: GetValidIngredientMeasurementUnit :one
SELECT
    valid_ingredient_measurement_units.id,
    valid_ingredient_measurement_units.notes,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on,
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
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on,
    valid_ingredient_measurement_units.created_on,
    valid_ingredient_measurement_units.last_updated_on,
    valid_ingredient_measurement_units.archived_on
FROM valid_ingredient_measurement_units
         JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
         JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_on IS NULL
  AND valid_ingredient_measurement_units.id = $1;

-- name: GetTotalValidIngredientMeasurementUnitsCount :one
SELECT COUNT(valid_ingredient_measurement_units.id) FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL;

-- name: CreateValidIngredientMeasurementUnit :exec
INSERT INTO valid_ingredient_measurement_units (id,notes,valid_measurement_unit_id,valid_ingredient_id) VALUES ($1,$2,$3,$4);

-- name: UpdateValidIngredientMeasurementUnit :exec
UPDATE valid_ingredient_measurement_units SET notes = $1, valid_measurement_unit_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4;

-- name: ArchiveValidIngredientMeasurementUnit :exec
UPDATE valid_ingredient_measurement_units SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;
