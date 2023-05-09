UPDATE recipe_steps SET
	index = $1,
	preparation_id = $2,
	time_scale_factor = $3,
	minimum_estimated_time_in_seconds = $4,
	maximum_estimated_time_in_seconds = $5,
	minimum_temperature_in_celsius = $6,
	maximum_temperature_in_celsius = $7,
	notes = $8,
	explicit_instructions = $9,
	condition_expression = $10,
	optional = $11,
	start_timer_automatically = $12,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe = $13
	AND id = $14;
