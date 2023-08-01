UPDATE recipe_step_completion_conditions
SET
	optional = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	ingredient_state = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $5;
