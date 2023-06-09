SELECT valid_ingredient_states.id,
       valid_ingredients.name,
       valid_ingredients.description,
       valid_ingredients.icon_path,
       valid_ingredients.slug,
       valid_ingredients.past_tense,
       valid_ingredients.attribute_type,
       valid_ingredient_states.created_at,
       valid_ingredient_states.last_updated_at,
       valid_ingredient_states.archived_at
  FROM valid_ingredient_states
 WHERE valid_ingredient_states.archived_at IS NULL;