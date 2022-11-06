-- name: ValidPreparationInstrumentExists :exec
SELECT EXISTS ( SELECT valid_preparation_instruments.id FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_at IS NULL AND valid_preparation_instruments.id = $1 );
