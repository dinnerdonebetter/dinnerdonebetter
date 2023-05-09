UPDATE recipe_step_vessels SET
	name = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	recipe_step_product_id = $4,
	valid_instrument_id = $5,
	vessel_predicate = $6,
	minimum_quantity = $7,
    maximum_quantity = $8,
    quantity_scale_factor = $9,
    unavailable_after_step = $10,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $11
	AND id = $12;