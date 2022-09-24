UPDATE valid_preparations
SET
	name = $1,
	description = $2,
	icon_path = $3,
	yields_nothing = $4,
	restrict_to_ingredients = $5,
	zero_ingredients_allowable = $6,
	past_tense = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
  AND id = $8;
