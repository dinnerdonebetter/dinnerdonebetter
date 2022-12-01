SELECT
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
    valid_preparations.minimum_ingredient_count,
    valid_preparations.maximum_ingredient_count,
    valid_preparations.minimum_instrument_count,
    valid_preparations.maximum_instrument_count,
    valid_preparations.temperature_required,
    valid_preparations.time_estimate_required,
    valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at
FROM valid_preparations
WHERE valid_preparations.archived_at IS NULL
	ORDER BY random() LIMIT 1;
