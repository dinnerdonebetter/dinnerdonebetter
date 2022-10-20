SELECT
	recipe_media.id,
    recipe_media.belongs_to_recipe,
    recipe_media.belongs_to_recipe_step,
    recipe_media.mime_type,
    recipe_media.internal_path,
    recipe_media.external_path,
	recipe_media.created_at,
	recipe_media.last_updated_at,
	recipe_media.archived_at
FROM recipe_media
WHERE recipe_media.archived_at IS NULL
	AND recipe_media.id = $1;
