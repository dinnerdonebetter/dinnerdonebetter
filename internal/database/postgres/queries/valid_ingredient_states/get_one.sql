SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
    valid_ingredient_states.slug,
    valid_ingredient_states.past_tense,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.id = $1;
