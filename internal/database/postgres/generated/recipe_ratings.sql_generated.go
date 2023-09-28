// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: recipe_ratings.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipeRating = `-- name: ArchiveRecipeRating :execrows

UPDATE recipe_ratings SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveRecipeRating(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveRecipeRating, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkRecipeRatingExistence = `-- name: CheckRecipeRatingExistence :one

SELECT EXISTS ( SELECT recipe_ratings.id FROM recipe_ratings WHERE recipe_ratings.archived_at IS NULL AND recipe_ratings.id = $1 )
`

func (q *Queries) CheckRecipeRatingExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeRatingExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipeRating = `-- name: CreateRecipeRating :exec

INSERT INTO recipe_ratings (
	id,
	recipe_id,
	taste,
	difficulty,
	cleanup,
	instructions,
	overall,
	notes,
	by_user
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9
)
`

type CreateRecipeRatingParams struct {
	ID           string
	RecipeID     string
	Notes        string
	ByUser       string
	Taste        sql.NullString
	Difficulty   sql.NullString
	Cleanup      sql.NullString
	Instructions sql.NullString
	Overall      sql.NullString
}

func (q *Queries) CreateRecipeRating(ctx context.Context, db DBTX, arg *CreateRecipeRatingParams) error {
	_, err := db.ExecContext(ctx, createRecipeRating,
		arg.ID,
		arg.RecipeID,
		arg.Taste,
		arg.Difficulty,
		arg.Cleanup,
		arg.Instructions,
		arg.Overall,
		arg.Notes,
		arg.ByUser,
	)
	return err
}

const getRecipeRating = `-- name: GetRecipeRating :one

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
	AND recipe_ratings.id = $1
`

func (q *Queries) GetRecipeRating(ctx context.Context, db DBTX, id string) (*RecipeRatings, error) {
	row := db.QueryRowContext(ctx, getRecipeRating, id)
	var i RecipeRatings
	err := row.Scan(
		&i.ID,
		&i.RecipeID,
		&i.Taste,
		&i.Difficulty,
		&i.Cleanup,
		&i.Instructions,
		&i.Overall,
		&i.Notes,
		&i.ByUser,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getRecipeRatings = `-- name: GetRecipeRatings :many

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
	 SELECT COUNT(recipe_ratings.id)
	 FROM recipe_ratings
	 WHERE recipe_ratings.archived_at IS NULL
	) as total_count
FROM
	recipe_ratings
WHERE
	recipe_ratings.archived_at IS NULL
	AND recipe_ratings.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND recipe_ratings.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
	AND (recipe_ratings.last_updated_at IS NULL OR recipe_ratings.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
GROUP BY recipe_ratings.id
ORDER BY recipe_ratings.id
OFFSET $5
LIMIT $6
`

type GetRecipeRatingsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipeRatingsRow struct {
	CreatedAt     time.Time
	ArchivedAt    sql.NullTime
	LastUpdatedAt sql.NullTime
	Notes         string
	RecipeID      string
	ID            string
	ByUser        string
	Difficulty    sql.NullString
	Overall       sql.NullString
	Instructions  sql.NullString
	Cleanup       sql.NullString
	Taste         sql.NullString
	FilteredCount int64
	TotalCount    int64
}

func (q *Queries) GetRecipeRatings(ctx context.Context, db DBTX, arg *GetRecipeRatingsParams) ([]*GetRecipeRatingsRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeRatings,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedAfter,
		arg.UpdatedBefore,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeRatingsRow{}
	for rows.Next() {
		var i GetRecipeRatingsRow
		if err := rows.Scan(
			&i.ID,
			&i.RecipeID,
			&i.Taste,
			&i.Difficulty,
			&i.Cleanup,
			&i.Instructions,
			&i.Overall,
			&i.Notes,
			&i.ByUser,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.FilteredCount,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRecipeRating = `-- name: UpdateRecipeRating :execrows

UPDATE recipe_ratings SET
	recipe_id = $1,
	taste = $2,
	difficulty = $3,
	cleanup = $4,
	instructions = $5,
	overall = $6,
	notes = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $8
`

type UpdateRecipeRatingParams struct {
	RecipeID     string
	Notes        string
	ID           string
	Taste        sql.NullString
	Difficulty   sql.NullString
	Cleanup      sql.NullString
	Instructions sql.NullString
	Overall      sql.NullString
}

func (q *Queries) UpdateRecipeRating(ctx context.Context, db DBTX, arg *UpdateRecipeRatingParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipeRating,
		arg.RecipeID,
		arg.Taste,
		arg.Difficulty,
		arg.Cleanup,
		arg.Instructions,
		arg.Overall,
		arg.Notes,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
