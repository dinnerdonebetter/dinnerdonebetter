-- name: ArchiveRecipeRating :exec

UPDATE recipe_ratings SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;


-- name: CreateRecipeRating :exec

INSERT INTO recipe_ratings (id,recipe_id,taste,difficulty,cleanup,instructions,overall,notes,by_user) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);


-- name: CheckRecipeRatingExistence :one

SELECT EXISTS ( SELECT recipe_ratings.id FROM recipe_ratings WHERE recipe_ratings.archived_at IS NULL AND recipe_ratings.id = $1 );

-- name: GetRecipeRatings :many

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
    recipe_ratings.archived_at,
	(
	 SELECT
		COUNT(recipe_ratings.id)
	 FROM
		recipe_ratings
	 WHERE
		recipe_ratings.archived_at IS NULL
	 AND recipe_ratings.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	 AND recipe_ratings.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	 AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	 AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
	 SELECT
		COUNT(recipe_ratings.id)
	 FROM
		recipe_ratings
	 WHERE
		recipe_ratings.archived_at IS NULL
	) as total_count
FROM
	recipe_ratings
WHERE
	recipe_ratings.archived_at IS NULL
	AND recipe_ratings.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND recipe_ratings.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
GROUP BY
	recipe_ratings.id
ORDER BY
	recipe_ratings.id
	LIMIT $5;


-- name: GetRecipeRating :one

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


-- name: UpdateRecipeRating :exec

UPDATE recipe_ratings
SET
	recipe_id = $1,
    taste = $2,
    difficulty = $3,
    cleanup = $4,
    instructions = $5,
    overall = $6,
    notes = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $8;
