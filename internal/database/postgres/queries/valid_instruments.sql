-- name: ValidInstrumentExists :one
SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.id = $1 );

-- name: GetValidInstrument :one
SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.plural_name,
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
    valid_instruments.plural_name,
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
    valid_instruments.plural_name,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on
FROM valid_instruments
WHERE valid_instruments.archived_on IS NULL
  AND valid_instruments.name ILIKE $1
  LIMIT 50;

-- name: GetValidInstruments :many
SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.plural_name,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.created_on,
    valid_instruments.last_updated_on,
    valid_instruments.archived_on,
    (
        SELECT
            COUNT(valid_instruments.id)
        FROM
            valid_instruments
        WHERE
            valid_instruments.archived_on IS NULL
          AND valid_instruments.created_on > COALESCE(sqlc.narg('created_after'), 0)
          AND valid_instruments.created_on < COALESCE(sqlc.narg('created_before'), (SELECT ~(1::bigint<<63)))
          AND (valid_instruments.last_updated_on IS NULL OR valid_instruments.last_updated_on > COALESCE(sqlc.narg('updated_after'), 0))
          AND (valid_instruments.last_updated_on IS NULL OR valid_instruments.last_updated_on < COALESCE(sqlc.narg('updated_before'), (SELECT ~(1::bigint<<63))))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_instruments.id)
        FROM
            valid_instruments
        WHERE
            valid_instruments.archived_on IS NULL
    ) as total_count
FROM
    valid_instruments
WHERE
    valid_instruments.archived_on IS NULL
  AND valid_instruments.created_on > COALESCE(sqlc.narg('created_after'), 0)
  AND valid_instruments.created_on < COALESCE(sqlc.narg('created_before'), (SELECT ~(1::bigint<<63)))
  AND (valid_instruments.last_updated_on IS NULL OR valid_instruments.last_updated_on > COALESCE(sqlc.narg('updated_after'), 0))
  AND (valid_instruments.last_updated_on IS NULL OR valid_instruments.last_updated_on < COALESCE(sqlc.narg('updated_before'), (SELECT ~(1::bigint<<63))))
GROUP BY
    valid_instruments.id
ORDER BY
    valid_instruments.id
    LIMIT sqlc.narg('limit');

-- name: CreateValidInstrument :exec
INSERT INTO valid_instruments (id,name,plural_name,description,icon_path) VALUES ($1,$2,$3,$4,$5);

-- name: UpdateValidInstrument :exec
UPDATE valid_instruments SET name = $1, plural_name = $2, description = $3, icon_path = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $5;

-- name: ArchiveValidInstrument :exec
UPDATE valid_instruments SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;