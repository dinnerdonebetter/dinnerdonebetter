UPDATE valid_ingredient_states
SET
	name = $1,
	description = $2,
	icon_path = $3,
    slug = $4,
	past_tense = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;