-- name: ArchiveRecipeStepVessel :exec

UPDATE recipe_step_vessels SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;
