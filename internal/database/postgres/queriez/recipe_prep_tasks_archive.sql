-- name: ArchiveRecipePrepTask :exec
UPDATE recipe_prep_tasks SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
