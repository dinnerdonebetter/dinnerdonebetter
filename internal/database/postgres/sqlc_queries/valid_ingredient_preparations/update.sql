-- name: UpdateValidIngredientPreparation :exec

UPDATE valid_ingredient_preparations SET notes = $1, valid_preparation_id = $2, valid_ingredient_id = $3, last_updated_at = NOW() WHERE archived_at IS NULL AND id = $4;
