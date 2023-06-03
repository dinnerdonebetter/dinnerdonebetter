-- name: GetMealRatings :many
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
    meal_ratings.archived_at,
	(
	 SELECT
		COUNT(meal_ratings.id)
	 FROM
		meal_ratings
	 WHERE
		meal_ratings.archived_at IS NULL
	 AND meal_ratings.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	 AND meal_ratings.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	 AND (meal_ratings.last_updated_at IS NULL OR meal_ratings.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	 AND (meal_ratings.last_updated_at IS NULL OR meal_ratings.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
	 SELECT
		COUNT(meal_ratings.id)
	 FROM
		meal_ratings
	 WHERE
		meal_ratings.archived_at IS NULL
	) as total_count
FROM
	meal_ratings
WHERE
	meal_ratings.archived_at IS NULL
	AND meal_ratings.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND meal_ratings.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	AND (meal_ratings.last_updated_at IS NULL OR meal_ratings.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	AND (meal_ratings.last_updated_at IS NULL OR meal_ratings.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
GROUP BY
	meal_ratings.id
ORDER BY
	meal_ratings.id
	LIMIT $5;
