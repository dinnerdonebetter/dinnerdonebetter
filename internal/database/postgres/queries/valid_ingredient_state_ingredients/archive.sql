-- name: ArchiveValidIngredientStateIngredient :exec

UPDATE valid_ingredient_state_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
