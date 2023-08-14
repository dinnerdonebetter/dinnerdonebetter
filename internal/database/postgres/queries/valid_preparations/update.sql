-- name: UpdateValidPreparation :exec

UPDATE valid_preparations
SET
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	icon_path = sqlc.arg(icon_path),
	yields_nothing = sqlc.arg(yields_nothing),
	restrict_to_ingredients = sqlc.arg(restrict_to_ingredients),
	minimum_ingredient_count = sqlc.arg(minimum_ingredient_count),
	maximum_ingredient_count = sqlc.arg(maximum_ingredient_count),
	minimum_instrument_count = sqlc.arg(minimum_instrument_count),
	maximum_instrument_count = sqlc.arg(maximum_instrument_count),
	temperature_required = sqlc.arg(temperature_required),
	time_estimate_required = sqlc.arg(time_estimate_required),
	condition_expression_required = sqlc.arg(condition_expression_required),
    consumes_vessel = sqlc.arg(consumes_vessel),
    only_for_vessels = sqlc.arg(only_for_vessels),
    minimum_vessel_count = sqlc.arg(minimum_vessel_count),
    maximum_vessel_count = sqlc.arg(maximum_vessel_count),
	slug = sqlc.arg(slug),
	past_tense = sqlc.arg(past_tense),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
