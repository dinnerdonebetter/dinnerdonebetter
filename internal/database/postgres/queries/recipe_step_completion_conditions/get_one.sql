-- name: GetRecipeStepCompletionCondition :many

SELECT
	recipe_step_completion_condition_ingredients.id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at,
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
	AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_step_completion_conditions.id = sqlc.arg(recipe_step_completion_condition_id)
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
	AND recipe_steps.id = sqlc.arg(recipe_step_id)
	AND recipes.archived_at IS NULL
	AND recipes.id = sqlc.arg(recipe_id);
