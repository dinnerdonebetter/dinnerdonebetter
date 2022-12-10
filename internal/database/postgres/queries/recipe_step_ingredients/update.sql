UPDATE recipe_step_ingredients SET
	ingredient_id = $1,
	name = $2,
	optional = $3,
	measurement_unit = $4,
	minimum_quantity_value = $5,
	maximum_quantity_value = $6,
	quantity_notes = $7,
	recipe_step_product_id = $8,
	ingredient_notes = $9,
	option_index = $10,
	requires_defrost = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND belongs_to_recipe_step = $12
	AND id = $13;
