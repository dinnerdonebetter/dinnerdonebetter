-- name: UpdateValidMeasurementConversion :exec

UPDATE valid_measurement_conversions
SET
	from_unit = $1,
	to_unit = $2,
	only_for_ingredient = $3,
	modifier = $4,
	notes = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
