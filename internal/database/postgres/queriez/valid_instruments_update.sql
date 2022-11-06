-- name: UpdateValidInstrument :exec
UPDATE valid_instruments
SET
	name = $1,
	plural_name = $2,
	description = $3,
	icon_path = $4,
	usable_for_storage = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
