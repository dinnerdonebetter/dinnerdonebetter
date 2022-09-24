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
