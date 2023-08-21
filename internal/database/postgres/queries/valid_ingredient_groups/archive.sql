-- name: ArchiveValidIngredientGroup :exec

UPDATE valid_ingredient_groups SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
