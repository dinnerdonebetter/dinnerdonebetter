// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: valid_ingredient_states.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const archiveValidIngredientState = `-- name: ArchiveValidIngredientState :execrows
UPDATE valid_ingredient_states SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveValidIngredientState(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveValidIngredientState, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkValidIngredientStateExistence = `-- name: CheckValidIngredientStateExistence :one
SELECT EXISTS (
	SELECT valid_ingredient_states.id
	FROM valid_ingredient_states
	WHERE valid_ingredient_states.archived_at IS NULL
		AND valid_ingredient_states.id = $1
)
`

func (q *Queries) CheckValidIngredientStateExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientStateExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createValidIngredientState = `-- name: CreateValidIngredientState :exec
INSERT INTO valid_ingredient_states (
	id,
	name,
	past_tense,
	slug,
	description,
	icon_path,
	attribute_type
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

type CreateValidIngredientStateParams struct {
	ID            string
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
}

func (q *Queries) CreateValidIngredientState(ctx context.Context, db DBTX, arg *CreateValidIngredientStateParams) error {
	_, err := db.ExecContext(ctx, createValidIngredientState,
		arg.ID,
		arg.Name,
		arg.PastTense,
		arg.Slug,
		arg.Description,
		arg.IconPath,
		arg.AttributeType,
	)
	return err
}

const getValidIngredientState = `-- name: GetValidIngredientState :one
SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.slug,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.last_indexed_at,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
AND valid_ingredient_states.id = $1
`

type GetValidIngredientStateRow struct {
	ID            string
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidIngredientState(ctx context.Context, db DBTX, id string) (*GetValidIngredientStateRow, error) {
	row := db.QueryRowContext(ctx, getValidIngredientState, id)
	var i GetValidIngredientStateRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PastTense,
		&i.Slug,
		&i.Description,
		&i.IconPath,
		&i.AttributeType,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidIngredientStates = `-- name: GetValidIngredientStates :many
SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.slug,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.last_indexed_at,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at,
	(
		SELECT COUNT(valid_ingredient_states.id)
		FROM valid_ingredient_states
		WHERE valid_ingredient_states.archived_at IS NULL
			AND
			valid_ingredient_states.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_ingredient_states.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_ingredient_states.last_updated_at IS NULL
				OR valid_ingredient_states.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_ingredient_states.last_updated_at IS NULL
				OR valid_ingredient_states.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE($5, false)::boolean OR valid_ingredient_states.archived_at = NULL)
	) AS filtered_count,
	(
		SELECT COUNT(valid_ingredient_states.id)
		FROM valid_ingredient_states
		WHERE valid_ingredient_states.archived_at IS NULL
	) AS total_count
FROM valid_ingredient_states
WHERE
	valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_ingredient_states.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_ingredient_states.last_updated_at IS NULL
		OR valid_ingredient_states.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_ingredient_states.last_updated_at IS NULL
		OR valid_ingredient_states.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
			AND (NOT COALESCE($5, false)::boolean OR valid_ingredient_states.archived_at = NULL)
GROUP BY valid_ingredient_states.id
ORDER BY valid_ingredient_states.id
LIMIT $7
OFFSET $6
`

type GetValidIngredientStatesParams struct {
	CreatedAfter    sql.NullTime
	CreatedBefore   sql.NullTime
	UpdatedBefore   sql.NullTime
	UpdatedAfter    sql.NullTime
	IncludeArchived sql.NullBool
	QueryOffset     sql.NullInt32
	QueryLimit      sql.NullInt32
}

type GetValidIngredientStatesRow struct {
	ID            string
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
	FilteredCount int64
	TotalCount    int64
}

func (q *Queries) GetValidIngredientStates(ctx context.Context, db DBTX, arg *GetValidIngredientStatesParams) ([]*GetValidIngredientStatesRow, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientStates,
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
	items := []*GetValidIngredientStatesRow{}
	for rows.Next() {
		var i GetValidIngredientStatesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PastTense,
			&i.Slug,
			&i.Description,
			&i.IconPath,
			&i.AttributeType,
			&i.LastIndexedAt,
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

const getValidIngredientStatesNeedingIndexing = `-- name: GetValidIngredientStatesNeedingIndexing :many
SELECT valid_ingredient_states.id
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND (
	valid_ingredient_states.last_indexed_at IS NULL
	OR valid_ingredient_states.last_indexed_at < NOW() - '24 hours'::INTERVAL
)
`

func (q *Queries) GetValidIngredientStatesNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientStatesNeedingIndexing)
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

const getValidIngredientStatesWithIDs = `-- name: GetValidIngredientStatesWithIDs :many
SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.slug,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.last_indexed_at,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.id = ANY($1::text[])
`

type GetValidIngredientStatesWithIDsRow struct {
	ID            string
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidIngredientStatesWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidIngredientStatesWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientStatesWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidIngredientStatesWithIDsRow{}
	for rows.Next() {
		var i GetValidIngredientStatesWithIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PastTense,
			&i.Slug,
			&i.Description,
			&i.IconPath,
			&i.AttributeType,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
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

const searchForValidIngredientStates = `-- name: SearchForValidIngredientStates :many
SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.slug,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.last_indexed_at,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.name ILIKE '%' || $1::text || '%'
	AND valid_ingredient_states.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientStatesRow struct {
	ID            string
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) SearchForValidIngredientStates(ctx context.Context, db DBTX, nameQuery string) ([]*SearchForValidIngredientStatesRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidIngredientStates, nameQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientStatesRow{}
	for rows.Next() {
		var i SearchForValidIngredientStatesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PastTense,
			&i.Slug,
			&i.Description,
			&i.IconPath,
			&i.AttributeType,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
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

const updateValidIngredientState = `-- name: UpdateValidIngredientState :execrows
UPDATE valid_ingredient_states SET
	name = $1,
	past_tense = $2,
	slug = $3,
	description = $4,
	icon_path = $5,
	attribute_type = $6,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $7
`

type UpdateValidIngredientStateParams struct {
	Name          string
	PastTense     string
	Slug          string
	Description   string
	IconPath      string
	AttributeType IngredientAttributeType
	ID            string
}

func (q *Queries) UpdateValidIngredientState(ctx context.Context, db DBTX, arg *UpdateValidIngredientStateParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidIngredientState,
		arg.Name,
		arg.PastTense,
		arg.Slug,
		arg.Description,
		arg.IconPath,
		arg.AttributeType,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateValidIngredientStateLastIndexedAt = `-- name: UpdateValidIngredientStateLastIndexedAt :execrows
UPDATE valid_ingredient_states SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateValidIngredientStateLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidIngredientStateLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
