UPDATE recipe_step_products
SET
	"name" = $1,
	"type" = $2,
	measurement_unit = $3,
	minimum_quantity_value = $4,
	maximum_quantity_value = $5,
	quantity_scale_factor = $6,
	quantity_notes = $7,
	compostable = $8,
	maximum_storage_duration_in_seconds = $9,
	minimum_storage_temperature_in_celsius = $10,
	maximum_storage_temperature_in_celsius = $11,
	storage_instructions = $12,
	is_liquid = $13,
	is_waste = $14,
    "index" = $15,
    contained_in_vessel_index = $16,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $17
	AND id = $18;
