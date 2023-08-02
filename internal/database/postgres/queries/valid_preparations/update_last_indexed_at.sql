-- name: UpdateValidPreparationLastIndexedAt :exec

UPDATE valid_preparations SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
