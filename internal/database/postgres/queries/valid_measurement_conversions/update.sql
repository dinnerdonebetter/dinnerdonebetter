UPDATE valid_measurement_conversions
SET
    from = $1
    to = $2
    only_for_ingredient = $3
    modifier = $4
    notes = $5
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
