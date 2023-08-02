-- name: UpdateRecipePrepTask :exec

UPDATE recipe_prep_tasks SET
	 name = $1,
	 description = $2,
	 notes = $3,
	 optional = $4,
	 explicit_storage_instructions = $5,
	 minimum_time_buffer_before_recipe_in_seconds = $6,
	 maximum_time_buffer_before_recipe_in_seconds = $7,
	 storage_type = $8,
	 minimum_storage_temperature_in_celsius = $9,
	 maximum_storage_temperature_in_celsius = $10,
	 belongs_to_recipe = $11,
	 last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $12;
