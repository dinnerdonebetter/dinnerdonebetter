-- name: CreateValidMeasurementUnitConversion :exec

INSERT INTO valid_measurement_conversions (id,from_unit,to_unit,only_for_ingredient,modifier,notes)
VALUES (sqlc.arg(id),sqlc.arg(from_unit),sqlc.arg(to_unit),sqlc.arg(only_for_ingredient),sqlc.arg(modifier)::float,sqlc.arg(notes));

