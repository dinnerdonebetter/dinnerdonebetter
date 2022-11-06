-- name: ArchiveRecipe :exec
UPDATE recipes SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2;
