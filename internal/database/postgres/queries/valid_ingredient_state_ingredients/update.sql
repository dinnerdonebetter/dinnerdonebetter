UPDATE valid_ingredient_state_ingredients SET notes = $1, valid_ingredient_state = $2, valid_ingredient = $3, last_updated_at = NOW() WHERE archived_at IS NULL AND id = $4;
