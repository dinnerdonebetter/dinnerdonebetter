SELECT
	recipe_ratings.id,
    recipe_ratings.recipe_id,
    recipe_ratings.taste,
    recipe_ratings.difficulty,
    recipe_ratings.cleanup,
    recipe_ratings.instructions,
    recipe_ratings.overall,
    recipe_ratings.notes,
    recipe_ratings.by_user,
    recipe_ratings.created_at,
    recipe_ratings.last_updated_at,
    recipe_ratings.archived_at
FROM recipe_ratings
WHERE recipe_ratings.archived_at IS NULL
	AND recipe_ratings.id = $1;
