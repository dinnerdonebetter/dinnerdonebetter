UPDATE recipe_step_instruments SET
	instrument_id = $1,
	recipe_step_product_id = $2,
	name = $3,
	notes = $4,
	preference_rank = $5,
	optional = $6,
	option_index = $7,
	minimum_quantity = $8,
	maximum_quantity = $9,
	quantity_scale_factor = $10,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $11
	AND id = $12;