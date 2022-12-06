SELECT EXISTS(
    SELECT id
    FROM valid_ingredient_state_ingredients
    WHERE valid_ingredient_id = $1
    AND valid_ingredient_state_id = $2
    AND archived_at IS NULL
);