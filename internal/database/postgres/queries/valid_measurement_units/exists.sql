-- name: CheckValidMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_at IS NULL AND valid_measurement_units.id = $1 );
