UPDATE meal_plan_grocery_list_items
SET
	belongs_to_meal_plan = $1,
	valid_ingredient = $2,
	valid_measurement_unit = $3,
	minimum_quantity_needed = $4,
	maximum_quantity_needed = $5,
	quantity_purchased = $6,
	purchased_measurement_unit = $7,
	purchased_upc = $8,
	purchase_price = $9,
	status_explanation = $10,
	status = $11,
	last_updated_at = NOW()
WHERE completed_at IS NULL
	AND id = $12;
