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
