-- name: UpdateValidInstrumentLastIndexedAt :exec

UPDATE valid_instruments SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
