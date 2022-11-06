-- name: ArchiveRecipeStepIngredient :exec
UPDATE recipe_step_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;
