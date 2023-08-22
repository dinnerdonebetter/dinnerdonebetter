-- name: ArchiveRecipeMedia :exec

UPDATE recipe_media SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateRecipeMedia :exec

INSERT INTO recipe_media (id,belongs_to_recipe,belongs_to_recipe_step,mime_type,internal_path,external_path,"index")
	VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: CheckRecipeMediaExistence :one

SELECT EXISTS ( SELECT recipe_media.id FROM recipe_media WHERE recipe_media.archived_at IS NULL AND recipe_media.id = $1 );

-- name: GetRecipeMediaForRecipe :many

SELECT
	recipe_media.id,
	recipe_media.belongs_to_recipe,
	recipe_media.belongs_to_recipe_step,
	recipe_media.mime_type,
	recipe_media.internal_path,
	recipe_media.external_path,
	recipe_media.index,
	recipe_media.created_at,
	recipe_media.last_updated_at,
	recipe_media.archived_at
FROM recipe_media
WHERE recipe_media.belongs_to_recipe = $1
	AND recipe_media.belongs_to_recipe_step IS NULL
	AND recipe_media.archived_at IS NULL
GROUP BY recipe_media.id
ORDER BY recipe_media.id;

-- name: GetRecipeMediaForRecipeStep :many

SELECT
	recipe_media.id,
	recipe_media.belongs_to_recipe,
	recipe_media.belongs_to_recipe_step,
	recipe_media.mime_type,
	recipe_media.internal_path,
	recipe_media.external_path,
	recipe_media.index,
	recipe_media.created_at,
	recipe_media.last_updated_at,
	recipe_media.archived_at
FROM recipe_media
WHERE recipe_media.belongs_to_recipe = $1
	AND recipe_media.belongs_to_recipe_step = $2
	AND recipe_media.archived_at IS NULL
GROUP BY recipe_media.id
ORDER BY recipe_media.id;

-- name: GetRecipeMedia :one

SELECT
	recipe_media.id,
	recipe_media.belongs_to_recipe,
	recipe_media.belongs_to_recipe_step,
	recipe_media.mime_type,
	recipe_media.internal_path,
	recipe_media.external_path,
	recipe_media.index,
	recipe_media.created_at,
	recipe_media.last_updated_at,
	recipe_media.archived_at
FROM recipe_media
WHERE recipe_media.archived_at IS NULL
	AND recipe_media.id = $1;

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
