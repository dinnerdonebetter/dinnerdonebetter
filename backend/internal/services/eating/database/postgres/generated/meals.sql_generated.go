// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: meals.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveMeal = `-- name: ArchiveMeal :execrows
UPDATE meals SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2
`

type ArchiveMealParams struct {
	CreatedByUser string
	ID            string
}

func (q *Queries) ArchiveMeal(ctx context.Context, db DBTX, arg *ArchiveMealParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveMeal, arg.CreatedByUser, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkMealExistence = `-- name: CheckMealExistence :one
SELECT EXISTS (
	SELECT meals.id
	FROM meals
	WHERE meals.archived_at IS NULL
		AND meals.id = $1
)
`

func (q *Queries) CheckMealExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createMeal = `-- name: CreateMeal :exec
INSERT INTO meals (
	id,
	name,
	description,
	min_estimated_portions,
	max_estimated_portions,
	eligible_for_meal_plans,
	created_by_user
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
`

type CreateMealParams struct {
	ID                   string
	Name                 string
	Description          string
	MinEstimatedPortions string
	CreatedByUser        string
	MaxEstimatedPortions sql.NullString
	EligibleForMealPlans bool
}

func (q *Queries) CreateMeal(ctx context.Context, db DBTX, arg *CreateMealParams) error {
	_, err := db.ExecContext(ctx, createMeal,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.MinEstimatedPortions,
		arg.MaxEstimatedPortions,
		arg.EligibleForMealPlans,
		arg.CreatedByUser,
	)
	return err
}

const getMeal = `-- name: GetMeal :many
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.id as component_id,
	meal_components.meal_id as component_meal_id,
	meal_components.recipe_id as component_recipe_id,
	meal_components.meal_component_type as component_meal_component_type,
	meal_components.recipe_scale as component_recipe_scale,
	meal_components.created_at as component_created_at,
	meal_components.last_updated_at as component_last_updated_at,
	meal_components.archived_at as component_archived_at
FROM meals
	JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
  AND meal_components.archived_at IS NULL
  AND meals.id = $1
`

type GetMealRow struct {
	CreatedAt                  time.Time
	ComponentCreatedAt         time.Time
	ComponentArchivedAt        sql.NullTime
	ComponentLastUpdatedAt     sql.NullTime
	ArchivedAt                 sql.NullTime
	LastUpdatedAt              sql.NullTime
	LastIndexedAt              sql.NullTime
	ComponentMealID            string
	CreatedByUser              string
	ComponentID                string
	ID                         string
	ComponentRecipeID          string
	ComponentMealComponentType ComponentType
	ComponentRecipeScale       string
	MinEstimatedPortions       string
	Description                string
	Name                       string
	MaxEstimatedPortions       sql.NullString
	EligibleForMealPlans       bool
}

func (q *Queries) GetMeal(ctx context.Context, db DBTX, id string) ([]*GetMealRow, error) {
	rows, err := db.QueryContext(ctx, getMeal, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMealRow{}
	for rows.Next() {
		var i GetMealRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.EligibleForMealPlans,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
			&i.ComponentID,
			&i.ComponentMealID,
			&i.ComponentRecipeID,
			&i.ComponentMealComponentType,
			&i.ComponentRecipeScale,
			&i.ComponentCreatedAt,
			&i.ComponentLastUpdatedAt,
			&i.ComponentArchivedAt,
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

const getMeals = `-- name: GetMeals :many
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND
			meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
	) AS filtered_count,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
	) AS total_count
FROM meals
WHERE
	meals.archived_at IS NULL
	AND meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
LIMIT $7
OFFSET $6
`

type GetMealsParams struct {
	CreatedAfter    sql.NullTime
	CreatedBefore   sql.NullTime
	UpdatedBefore   sql.NullTime
	UpdatedAfter    sql.NullTime
	IncludeArchived sql.NullBool
	QueryOffset     sql.NullInt32
	QueryLimit      sql.NullInt32
}

type GetMealsRow struct {
	CreatedAt            time.Time
	ArchivedAt           sql.NullTime
	LastIndexedAt        sql.NullTime
	LastUpdatedAt        sql.NullTime
	Description          string
	MinEstimatedPortions string
	Name                 string
	ID                   string
	CreatedByUser        string
	MaxEstimatedPortions sql.NullString
	FilteredCount        int64
	TotalCount           int64
	EligibleForMealPlans bool
}

func (q *Queries) GetMeals(ctx context.Context, db DBTX, arg *GetMealsParams) ([]*GetMealsRow, error) {
	rows, err := db.QueryContext(ctx, getMeals,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.IncludeArchived,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMealsRow{}
	for rows.Next() {
		var i GetMealsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.EligibleForMealPlans,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
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

const getMealsCreatedByUser = `-- name: GetMealsCreatedByUser :many
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND
			meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
			AND meals.created_by_user = $6
	) AS filtered_count,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND meals.created_by_user = $6
	) AS total_count
FROM meals
WHERE
	meals.archived_at IS NULL
	AND meals.created_by_user = $6
	AND meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
	AND meals.created_by_user = $6
LIMIT $8
OFFSET $7
`

type GetMealsCreatedByUserParams struct {
	CreatedAfter    sql.NullTime
	CreatedBefore   sql.NullTime
	UpdatedBefore   sql.NullTime
	UpdatedAfter    sql.NullTime
	CreatedByUser   string
	QueryOffset     sql.NullInt32
	QueryLimit      sql.NullInt32
	IncludeArchived sql.NullBool
}

type GetMealsCreatedByUserRow struct {
	CreatedAt            time.Time
	ArchivedAt           sql.NullTime
	LastIndexedAt        sql.NullTime
	LastUpdatedAt        sql.NullTime
	Description          string
	MinEstimatedPortions string
	Name                 string
	ID                   string
	CreatedByUser        string
	MaxEstimatedPortions sql.NullString
	FilteredCount        int64
	TotalCount           int64
	EligibleForMealPlans bool
}

func (q *Queries) GetMealsCreatedByUser(ctx context.Context, db DBTX, arg *GetMealsCreatedByUserParams) ([]*GetMealsCreatedByUserRow, error) {
	rows, err := db.QueryContext(ctx, getMealsCreatedByUser,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.IncludeArchived,
		arg.CreatedByUser,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMealsCreatedByUserRow{}
	for rows.Next() {
		var i GetMealsCreatedByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.EligibleForMealPlans,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
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

const getMealsNeedingIndexing = `-- name: GetMealsNeedingIndexing :many
SELECT meals.id
	FROM meals
	WHERE meals.archived_at IS NULL
	AND (
		meals.last_indexed_at IS NULL
		OR meals.last_indexed_at < NOW() - '24 hours'::INTERVAL
	)
`

func (q *Queries) GetMealsNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getMealsNeedingIndexing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchForMeals = `-- name: SearchForMeals :many
SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.id as component_id,
	meal_components.meal_id as component_meal_id,
	meal_components.recipe_id as component_recipe_id,
	meal_components.meal_component_type as component_meal_component_type,
	meal_components.recipe_scale as component_recipe_scale,
	meal_components.created_at as component_created_at,
	meal_components.last_updated_at as component_last_updated_at,
	meal_components.archived_at as component_archived_at,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND
			meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
	) AS filtered_count,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
	) AS total_count
FROM meals
	JOIN meal_components ON meal_components.meal_id=meals.id
WHERE
	meals.archived_at IS NULL
	AND meals.name ILIKE '%' || $6::text || '%'
	AND meals.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND meals.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
			AND (NOT COALESCE($5, false)::boolean OR meals.archived_at = NULL)
LIMIT $8
OFFSET $7
`

type SearchForMealsParams struct {
	CreatedAfter    sql.NullTime
	CreatedBefore   sql.NullTime
	UpdatedBefore   sql.NullTime
	UpdatedAfter    sql.NullTime
	Query           string
	QueryOffset     sql.NullInt32
	QueryLimit      sql.NullInt32
	IncludeArchived sql.NullBool
}

type SearchForMealsRow struct {
	CreatedAt                  time.Time
	ComponentCreatedAt         time.Time
	ComponentArchivedAt        sql.NullTime
	ComponentLastUpdatedAt     sql.NullTime
	ArchivedAt                 sql.NullTime
	LastUpdatedAt              sql.NullTime
	LastIndexedAt              sql.NullTime
	ComponentRecipeScale       string
	ComponentMealComponentType ComponentType
	Name                       string
	CreatedByUser              string
	ComponentID                string
	ComponentMealID            string
	ComponentRecipeID          string
	Description                string
	ID                         string
	MinEstimatedPortions       string
	MaxEstimatedPortions       sql.NullString
	FilteredCount              int64
	TotalCount                 int64
	EligibleForMealPlans       bool
}

func (q *Queries) SearchForMeals(ctx context.Context, db DBTX, arg *SearchForMealsParams) ([]*SearchForMealsRow, error) {
	rows, err := db.QueryContext(ctx, searchForMeals,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.IncludeArchived,
		arg.Query,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForMealsRow{}
	for rows.Next() {
		var i SearchForMealsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.EligibleForMealPlans,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
			&i.ComponentID,
			&i.ComponentMealID,
			&i.ComponentRecipeID,
			&i.ComponentMealComponentType,
			&i.ComponentRecipeScale,
			&i.ComponentCreatedAt,
			&i.ComponentLastUpdatedAt,
			&i.ComponentArchivedAt,
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

const updateMealLastIndexedAt = `-- name: UpdateMealLastIndexedAt :execrows
UPDATE meals SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateMealLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateMealLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
