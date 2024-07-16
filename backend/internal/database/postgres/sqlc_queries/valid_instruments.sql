-- name: ArchiveValidInstrument :execrows

UPDATE valid_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND id = sqlc.arg(id);
-- name: CreateValidInstrument :exec

INSERT INTO valid_instruments (
	id,
	name,
	description,
	icon_path,
	plural_name,
	usable_for_storage,
	slug,
	display_in_summary_lists,
	include_in_generated_instructions
) VALUES (
	sqlc.arg(id),
	sqlc.arg(name),
	sqlc.arg(description),
	sqlc.arg(icon_path),
	sqlc.arg(plural_name),
	sqlc.arg(usable_for_storage),
	sqlc.arg(slug),
	sqlc.arg(display_in_summary_lists),
	sqlc.arg(include_in_generated_instructions)
);
-- name: CheckValidInstrumentExistence :one

SELECT EXISTS (
	SELECT valid_instruments.id
	FROM valid_instruments
	WHERE valid_instruments.archived_at IS NULL
		AND valid_instruments.id = sqlc.arg(id)
);
-- name: GetValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	(
		SELECT COUNT(valid_instruments.id)
		FROM valid_instruments
		WHERE valid_instruments.archived_at IS NULL
			AND valid_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_instruments.last_updated_at IS NULL
				OR valid_instruments.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_instruments.last_updated_at IS NULL
				OR valid_instruments.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(valid_instruments.id)
		FROM valid_instruments
		WHERE valid_instruments.archived_at IS NULL
	) AS total_count
FROM valid_instruments
WHERE
	valid_instruments.archived_at IS NULL
	AND valid_instruments.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_instruments.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_instruments.last_updated_at IS NULL
		OR valid_instruments.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_instruments.last_updated_at IS NULL
		OR valid_instruments.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
GROUP BY valid_instruments.id
ORDER BY valid_instruments.id
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
-- name: GetValidInstrumentsNeedingIndexing :many

SELECT valid_instruments.id
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND (
	valid_instruments.last_indexed_at IS NULL
	OR valid_instruments.last_indexed_at < NOW() - '24 hours'::INTERVAL
);
-- name: GetValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
AND valid_instruments.id = sqlc.arg(id);
-- name: GetRandomValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
ORDER BY RANDOM() LIMIT 1;
-- name: GetValidInstrumentsWithIDs :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
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
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.name ILIKE '%' || sqlc.arg(name_query)::text || '%'
	AND valid_instruments.archived_at IS NULL
LIMIT 50;
-- name: UpdateValidInstrument :execrows

UPDATE valid_instruments SET
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	icon_path = sqlc.arg(icon_path),
	plural_name = sqlc.arg(plural_name),
	usable_for_storage = sqlc.arg(usable_for_storage),
	slug = sqlc.arg(slug),
	display_in_summary_lists = sqlc.arg(display_in_summary_lists),
	include_in_generated_instructions = sqlc.arg(include_in_generated_instructions),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
-- name: UpdateValidInstrumentLastIndexedAt :execrows

UPDATE valid_instruments SET last_indexed_at = NOW() WHERE id = sqlc.arg(id) AND archived_at IS NULL;
