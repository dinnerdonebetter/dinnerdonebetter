-- name: GetValidVesselIDsNeedingIndexing :many

SELECT
	valid_vessels.id
FROM valid_vessels
	 JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE (valid_vessels.archived_at IS NULL AND valid_measurement_units.archived_at IS NULL)
   AND (
        (valid_vessels.last_indexed_at IS NULL)
        OR valid_vessels.last_indexed_at
            < now() - '24 hours'::INTERVAL
    );