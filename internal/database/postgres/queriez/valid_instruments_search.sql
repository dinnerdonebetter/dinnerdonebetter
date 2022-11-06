-- name: SearchForValidInstruments :exec
SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.name ILIKE $1 LIMIT 50;
