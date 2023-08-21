-- name: CheckValidIngredientMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_at IS NULL AND valid_ingredient_measurement_units.id = $1 );
