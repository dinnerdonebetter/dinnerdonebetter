-- name: ValidInstrumentExists :one
SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id = $1 );

-- name: GetValidInstrument :one
SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on
FROM valid_instruments
WHERE valid_instruments.archived_on IS NULL
AND valid_instruments.id = $1;

-- name: GetRandomValidInstrument :one
SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on
FROM valid_instruments
WHERE valid_instruments.archived_on IS NULL
ORDER BY random() LIMIT 1;

-- name: SearchForValidInstruments :many
SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.variant,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on
FROM valid_instruments
WHERE valid_instruments.archived_on IS NULL
  AND valid_instruments.name ILIKE $1
  LIMIT 50;

-- name: GetTotalValidInstrumentCount :one
SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL;

-- name: CreateValidInstrument :exec
INSERT INTO valid_instruments (id,name,variant,description,icon_path) VALUES ($1,$2,$3,$4,$5);

-- name: UpdateValidInstrument :exec
UPDATE valid_instruments SET name = $1, variant = $2, description = $3, icon_path = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $5;

-- name: ArchiveValidInstrument :exec
UPDATE valid_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;