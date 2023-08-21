-- name: ArchiveValidVessel :exec

UPDATE valid_vessels SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
