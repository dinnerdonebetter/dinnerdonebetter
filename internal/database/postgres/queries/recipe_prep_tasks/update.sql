UPDATE recipe_prep_tasks SET
     notes = $1,
     explicit_storage_instructions = $2,
     minimum_time_buffer_before_recipe_in_seconds = $3,
     maximum_time_buffer_before_recipe_in_seconds = $4,
     storage_type = $5,
     minimum_storage_temperature_in_celsius = $6,
     maximum_storage_temperature_in_celsius = $7,
     belongs_to_recipe = $8,
     last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $9;
