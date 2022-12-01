UPDATE valid_preparations
SET
	name = $1,
	description = $2,
	icon_path = $3,
	yields_nothing = $4,
	restrict_to_ingredients = $5,
	zero_ingredients_allowable = $6,
    minimum_ingredient_count = $7,
    maximum_ingredient_count = $8,
    minimum_instrument_count = $9,
    maximum_instrument_count = $10,
    temperature_required = $11,
    time_estimate_required = $12,
    slug = $13,
	past_tense = $14,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $15;
