-- name: ValidIngredientExists :exec
SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_at IS NULL AND valid_ingredients.id = $1 );
