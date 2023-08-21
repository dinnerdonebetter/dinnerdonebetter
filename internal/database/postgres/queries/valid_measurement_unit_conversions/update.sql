-- name: UpdateValidMeasurementUnitConversion :exec

UPDATE valid_measurement_conversions
SET
	from_unit = sqlc.arg(from_unit),
	to_unit = sqlc.arg(to_unit),
	only_for_ingredient = sqlc.arg(only_for_ingredient),
	modifier = sqlc.arg(modifier)::float,
	notes = sqlc.arg(notes),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
