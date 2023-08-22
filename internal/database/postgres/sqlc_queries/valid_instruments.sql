-- name: ArchiveValidInstrument :exec

UPDATE valid_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateValidInstrument :exec

INSERT INTO valid_instruments (id,"name",plural_name,description,icon_path,usable_for_storage,display_in_summary_lists,include_in_generated_instructions,slug) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: CheckValidInstrumentExistence :one

SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_at IS NULL AND valid_instruments.id = $1 );

-- name: GetValidInstruments :many

SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.plural_name,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.usable_for_storage,
    valid_instruments.display_in_summary_lists,
    valid_instruments.include_in_generated_instructions,
    valid_instruments.slug,
    valid_instruments.created_at,
    valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	(
	 SELECT
		COUNT(valid_instruments.id)
	 FROM
		valid_instruments
	 WHERE
		valid_instruments.archived_at IS NULL
	 AND valid_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	 AND valid_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	 AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	 AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
	 SELECT
		COUNT(valid_instruments.id)
	 FROM
		valid_instruments
	 WHERE
		valid_instruments.archived_at IS NULL
	) as total_count
FROM
	valid_instruments
WHERE
	valid_instruments.archived_at IS NULL
	AND valid_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND valid_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
GROUP BY
	valid_instruments.id
ORDER BY
	valid_instruments.id
    OFFSET sqlc.narg(query_offset)
	LIMIT sqlc.narg(query_limit);

-- name: GetValidInstrumentsNeedingIndexing :many

SELECT valid_instruments.id
  FROM valid_instruments
 WHERE (valid_instruments.archived_at IS NULL)
       AND (
			(valid_instruments.last_indexed_at IS NULL)
			OR valid_instruments.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.id = $1;

-- name: GetRandomValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	ORDER BY random() LIMIT 1;

-- name: GetValidInstrumentWithIDs :many

SELECT
    valid_instruments.id,
    valid_instruments.name,
    valid_instruments.plural_name,
    valid_instruments.description,
    valid_instruments.icon_path,
    valid_instruments.usable_for_storage,
    valid_instruments.display_in_summary_lists,
    valid_instruments.include_in_generated_instructions,
    valid_instruments.slug,
    valid_instruments.created_at,
    valid_instruments.last_updated_at,
    valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
    AND valid_instruments.id = ANY(sqlc.arg(ids)::text[]);

-- name: SearchForValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.name ILIKE '%' || sqlc.arg(query)::text || '%'
    LIMIT 50;

-- name: UpdateValidInstrument :exec

UPDATE valid_instruments
SET
	name = $1,
	plural_name = $2,
	description = $3,
	icon_path = $4,
	usable_for_storage = $5,
	display_in_summary_lists = $6,
	include_in_generated_instructions = $7,
	slug = $8,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $9;

-- name: UpdateValidInstrumentLastIndexedAt :exec

UPDATE valid_instruments SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
