-- name: CheckRecipeRatingExistence :one

SELECT EXISTS ( SELECT recipe_ratings.id FROM recipe_ratings WHERE recipe_ratings.archived_at IS NULL AND recipe_ratings.id = $1 );