SELECT
    recipe_step_completion_conditions.id,
    recipe_step_completion_conditions.belongs_to_recipe_step,
    recipe_step_completion_conditions.ingredient_state,
    recipe_step_completion_conditions.optional,
    recipe_step_completion_conditions.notes,
    recipe_step_completion_conditions.created_at,
    recipe_step_completion_conditions.last_updated_at,
    recipe_step_completion_conditions.archived_at
FROM recipe_step_completion_conditions
    JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
    JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
    JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE
    recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_conditions.belongs_to_recipe_step = $2
	AND recipe_step_completion_conditions.id = $3
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipe_steps.id = $2
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;
