-- name: GetValidMeasurementUnitsNeedingIndexing :many

SELECT valid_measurement_units.id
  FROM valid_measurement_units
 WHERE (valid_measurement_units.archived_at IS NULL)
       AND (
			(
				valid_measurement_units.last_indexed_at IS NULL
			)
			OR valid_measurement_units.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);