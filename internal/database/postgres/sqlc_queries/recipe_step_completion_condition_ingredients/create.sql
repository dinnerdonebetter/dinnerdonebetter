INSERT INTO recipe_step_completion_condition_ingredients (
	id,
	belongs_to_recipe_step_completion_condition,
	recipe_step_ingredient
) VALUES ($1,$2,$3);
