UPDATE recipe_step_ingredients SET
	ingredient_id = $1,
	name = $2,
	optional = $3,
	measurement_unit = $4,
	minimum_quantity_value = $5,
	maximum_quantity_value = $6,
	quantity_scale_factor = $7,
	quantity_notes = $8,
	recipe_step_product_id = $9,
	ingredient_notes = $10,
	option_index = $11,
	to_taste = $12,
	product_percentage_to_use = $13,
    vessel_index = $14,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND belongs_to_recipe_step = $15
	AND id = $16;
