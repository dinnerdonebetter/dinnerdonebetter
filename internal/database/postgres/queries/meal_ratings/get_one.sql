SELECT
	meal_ratings.id,
    meal_ratings.meal_id,
    meal_ratings.taste,
    meal_ratings.difficulty,
    meal_ratings.cleanup,
    meal_ratings.instructions,
    meal_ratings.overall,
    meal_ratings.notes,
    meal_ratings.by_user,
    meal_ratings.created_at,
    meal_ratings.last_updated_at,
    meal_ratings.archived_at
FROM meal_ratings
WHERE meal_ratings.archived_at IS NULL
	AND meal_ratings.id = $1;
