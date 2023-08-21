-- name: UpdateValidVesselLastIndexedAt :exec

UPDATE valid_vessels SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
