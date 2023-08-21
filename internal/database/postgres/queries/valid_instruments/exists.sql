-- name: CheckValidInstrumentExistence :one

SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_at IS NULL AND valid_instruments.id = $1 );