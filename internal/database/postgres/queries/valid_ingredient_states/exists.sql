-- name: CheckValidIngredientStateExistence :one

SELECT EXISTS ( SELECT valid_ingredient_states.id FROM valid_ingredient_states WHERE valid_ingredient_states.archived_at IS NULL AND valid_ingredient_states.id = $1 );
