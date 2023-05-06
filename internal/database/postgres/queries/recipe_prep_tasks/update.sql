UPDATE recipe_prep_tasks SET
	 name = $1,
	 description = $2,
	 notes = $3,
	 explicit_storage_instructions = $4,
	 minimum_time_buffer_before_recipe_in_seconds = $5,
	 maximum_time_buffer_before_recipe_in_seconds = $6,
	 storage_type = $7,
	 minimum_storage_temperature_in_celsius = $8,
	 maximum_storage_temperature_in_celsius = $9,
	 belongs_to_recipe = $10,
	 last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $11;
