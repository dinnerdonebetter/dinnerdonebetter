-- name: UpdateValidPreparation :exec

UPDATE valid_preparations
SET
	name = $1,
	description = $2,
	icon_path = $3,
	yields_nothing = $4,
	restrict_to_ingredients = $5,
	minimum_ingredient_count = $6,
	maximum_ingredient_count = $7,
	minimum_instrument_count = $8,
	maximum_instrument_count = $9,
	temperature_required = $10,
	time_estimate_required = $11,
	condition_expression_required = $12,
    consumes_vessel = $13,
    only_for_vessels = $14,
    minimum_vessel_count = $15,
    maximum_vessel_count = $16,
	slug = $17,
	past_tense = $18,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $19;
