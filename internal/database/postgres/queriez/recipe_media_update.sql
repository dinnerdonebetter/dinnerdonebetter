-- name: UpdateRecipeMedia :exec
UPDATE recipe_media
SET
    belongs_to_recipe = $1,
    belongs_to_recipe_step = $2,
    mime_type = $3,
    internal_path = $4,
    external_path = $5,
    "index" = $6,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
