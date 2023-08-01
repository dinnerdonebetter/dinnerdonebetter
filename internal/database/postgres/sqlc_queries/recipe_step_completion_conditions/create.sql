INSERT INTO recipe_step_completion_conditions (
	id,
	belongs_to_recipe_step,
	ingredient_state,
	optional,
	notes
) VALUES ($1,$2,$3,$4,$5);
