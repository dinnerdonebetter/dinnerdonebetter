-- name: ArchiveValidInstrument :exec
UPDATE valid_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;