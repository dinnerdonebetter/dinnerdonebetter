-- name: ValidPreparationInstrumentExists :one
SELECT EXISTS ( SELECT valid_preparation_instruments.id FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL AND valid_preparation_instruments.id = $1 );

-- name: GetValidPreparationInstrument :one
SELECT
    valid_preparation_instruments.id,
    valid_preparation_instruments.notes,
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on,
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on,
    valid_preparation_instruments.created_on,
    valid_preparation_instruments.last_updated_on,
    valid_preparation_instruments.archived_on
FROM
    valid_preparation_instruments
        JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
        JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
    valid_preparation_instruments.archived_on IS NULL
  AND valid_preparation_instruments.id = $1;

-- name: GetTotalValidPreparationInstrumentsCount :one
SELECT COUNT(valid_preparation_instruments.id) FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_on IS NULL;

-- name: CreateValidPreparationInstrument :exec
INSERT INTO valid_preparation_instruments (id,notes,valid_preparation_id,valid_instrument_id) VALUES ($1,$2,$3,$4);

-- name: UpdateValidPreparationInstrument :exec
UPDATE valid_preparation_instruments SET notes = $1, valid_preparation_id = $2, valid_instrument_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4;

-- name: ArchiveValidPreparationInstrument :exec
UPDATE valid_preparation_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;
