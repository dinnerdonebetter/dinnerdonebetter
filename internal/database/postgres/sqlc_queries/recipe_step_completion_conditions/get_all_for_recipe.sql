-- name: GetAllRecipeStepCompletionConditionsForRecipe :many

SELECT
	recipe_step_completion_condition_ingredients.id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.slug,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at,
	recipe_step_completion_conditions.optional,
	recipe_step_completion_conditions.notes,
	recipe_step_completion_conditions.created_at,
	recipe_step_completion_conditions.last_updated_at,
	recipe_step_completion_conditions.archived_at
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipes.archived_at IS NULL
	AND recipes.id = $1
GROUP BY recipe_step_completion_conditions.id,
	     recipe_step_completion_condition_ingredients.id,
	     valid_ingredient_states.id;