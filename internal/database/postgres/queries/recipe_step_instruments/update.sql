UPDATE recipe_step_instruments SET
	instrument_id = $1,
	recipe_step_product_id = $2,
	name = $3,
	product_of_recipe_step = $4,
	notes = $5,
	preference_rank = $6,
	optional = $7,
    option_index = $8,
	minimum_quantity = $9,
    maximum_quantity = $10,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $11
	AND id = $12;