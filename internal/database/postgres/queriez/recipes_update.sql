-- name: UpdateRecipe :exec
UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, yields_portions = $5, seal_of_approval = $6, last_updated_at = NOW() WHERE archived_at IS NULL AND created_by_user = $7 AND id = $8;
