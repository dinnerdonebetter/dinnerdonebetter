-- name: CreateRecipeStepIngredient :exec

INSERT INTO recipe_step_ingredients (
	id,
	"name",
	optional,
	ingredient_id,
	measurement_unit,
	minimum_quantity_value,
	maximum_quantity_value,
	quantity_notes,
	recipe_step_product_id,
	ingredient_notes,
	option_index,
	to_taste,
	product_percentage_to_use,
    vessel_index,
    recipe_step_product_recipe_id,
	belongs_to_recipe_step
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);
