-- name: ArchiveValidPreparation :exec

UPDATE valid_preparations SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
