-- name: GetRecipeStepCompletionConditions :many

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
	recipe_step_completion_conditions.archived_at,
	(
	    SELECT
	        COUNT(recipe_step_completion_conditions.id)
	    FROM
	        recipe_step_completion_conditions
	    WHERE
	        recipe_step_completion_conditions.archived_at IS NULL
	      AND recipe_step_completion_conditions.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	      AND recipe_step_completion_conditions.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	      AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	      AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
	    SELECT
	        COUNT(recipe_step_completion_conditions.id)
	    FROM
	        recipe_step_completion_conditions
	    WHERE
	        recipe_step_completion_conditions.archived_at IS NULL
	) as total_count
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE
	recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_conditions.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND recipe_step_completion_conditions.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
GROUP BY
	recipe_step_completion_conditions.id,
	valid_ingredient_states.id,
	recipe_step_completion_condition_ingredients.id
ORDER BY
	recipe_step_completion_conditions.id
	LIMIT $5
	OFFSET $6;