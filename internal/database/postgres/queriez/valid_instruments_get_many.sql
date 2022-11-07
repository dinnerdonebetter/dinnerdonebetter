-- name: GetValidInstruments :many
SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
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
	 AND valid_instruments.created_at > COALESCE(sqlc.narg('created_after'), 0)
	 AND valid_instruments.created_at < COALESCE(sqlc.narg('created_before'), (SELECT NOW() + interval '999 years'))
	 AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at > COALESCE(sqlc.narg('updated_after'), 0))
	 AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at < COALESCE(sqlc.narg('updated_before'), (SELECT NOW() + interval '999 years')))
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
	AND valid_instruments.created_at > COALESCE(sqlc.narg('created_after'), 0)
	AND valid_instruments.created_at < COALESCE(sqlc.narg('created_before'), (SELECT NOW() + interval '999 years'))
	AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at > COALESCE(sqlc.narg('updated_after'), 0))
	AND (valid_instruments.last_updated_at IS NULL OR valid_instruments.last_updated_at < COALESCE(sqlc.narg('updated_before'), (SELECT NOW() + interval '999 years')))
GROUP BY
	valid_instruments.id
ORDER BY
	valid_instruments.id
	LIMIT sqlc.narg('limit');
