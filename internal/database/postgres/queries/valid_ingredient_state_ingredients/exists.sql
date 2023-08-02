-- name: CheckValidIngredientStateIngredientExistence :one

SELECT EXISTS ( SELECT valid_ingredient_state_ingredients.id FROM valid_ingredient_state_ingredients WHERE valid_ingredient_state_ingredients.archived_at IS NULL AND valid_ingredient_state_ingredients.id = $1 );
