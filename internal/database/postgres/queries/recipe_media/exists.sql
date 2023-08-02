-- name: CheckRecipeMediaExistence :one

SELECT EXISTS ( SELECT recipe_media.id FROM recipe_media WHERE recipe_media.archived_at IS NULL AND recipe_media.id = $1 );
