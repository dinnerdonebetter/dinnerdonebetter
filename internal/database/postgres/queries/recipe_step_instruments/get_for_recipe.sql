SELECT
	recipe_step_instruments.id,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
    valid_instruments.is_vessel,
    valid_instruments.is_exclusively_vessel,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.optional,
	recipe_step_instruments.minimum_quantity,
	recipe_step_instruments.maximum_quantity,
	recipe_step_instruments.option_index,
	recipe_step_instruments.created_at,
	recipe_step_instruments.last_updated_at,
	recipe_step_instruments.archived_at,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
	LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
	JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1;
