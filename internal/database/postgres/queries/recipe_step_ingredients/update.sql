UPDATE recipe_step_ingredients SET
	ingredient_id = $1,
	name = $2,
	optional = $3,
	measurement_unit = $4,
	minimum_quantity_value = $5,
	maximum_quantity_value = $6,
	quantity_notes = $7,
	product_of_recipe_step = $8,
	recipe_step_product_id = $9,
    ingredient_notes = $10,
    option_index = $11,
    requires_defrost = $12,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND belongs_to_recipe_step = $13
	AND id = $14;
