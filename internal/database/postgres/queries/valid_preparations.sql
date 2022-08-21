-- name: ValidPreparationExists :one
SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.id = $1 );

-- name: GetValidPreparation :one
SELECT
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on
FROM valid_preparations
WHERE valid_preparations.archived_on IS NULL
  AND valid_preparations.id = $1;

-- name: GetRandomValidPreparation :one
SELECT
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on
FROM valid_preparations
WHERE valid_preparations.archived_on IS NULL
    ORDER BY random() LIMIT 1;

-- name: ValidPreparationsSearch :many
SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.name ILIKE $1 LIMIT 50;

-- name: GetTotalValidPreparationsCount :one
SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL;

-- name: CreateValidPreparation :exec
INSERT INTO valid_preparations (id,name,description,icon_path) VALUES ($1,$2,$3,$4);

-- name: UpdateValidPreparation :exec
UPDATE valid_preparations SET name = $1, description = $2, icon_path = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $4;

-- name: ArchiveValidPreparation :exec
UPDATE valid_preparations SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;
