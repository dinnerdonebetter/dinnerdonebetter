UPDATE valid_ingredient_groups SET
	name = $1,
	description = $2,
	slug = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $4;
