// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: valid_measurement_units.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const archiveValidMeasurementUnit = `-- name: ArchiveValidMeasurementUnit :execrows
UPDATE valid_measurement_units SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveValidMeasurementUnit(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveValidMeasurementUnit, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkValidMeasurementUnitExistence = `-- name: CheckValidMeasurementUnitExistence :one
SELECT EXISTS (
	SELECT valid_measurement_units.id
	FROM valid_measurement_units
	WHERE valid_measurement_units.archived_at IS NULL
		AND valid_measurement_units.id = $1
)
`

func (q *Queries) CheckValidMeasurementUnitExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidMeasurementUnitExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createValidMeasurementUnit = `-- name: CreateValidMeasurementUnit :exec
INSERT INTO valid_measurement_units (
	id,
	name,
	description,
	volumetric,
	icon_path,
	universal,
	metric,
	imperial,
	slug,
	plural_name
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10
)
`

type CreateValidMeasurementUnitParams struct {
	ID          string
	Name        string
	Description string
	Volumetric  sql.NullBool
	IconPath    string
	Universal   bool
	Metric      bool
	Imperial    bool
	Slug        string
	PluralName  string
}

func (q *Queries) CreateValidMeasurementUnit(ctx context.Context, db DBTX, arg *CreateValidMeasurementUnitParams) error {
	_, err := db.ExecContext(ctx, createValidMeasurementUnit,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Volumetric,
		arg.IconPath,
		arg.Universal,
		arg.Metric,
		arg.Imperial,
		arg.Slug,
		arg.PluralName,
	)
	return err
}

const getRandomValidMeasurementUnit = `-- name: GetRandomValidMeasurementUnit :one
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
ORDER BY RANDOM() LIMIT 1
`

type GetRandomValidMeasurementUnitRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetRandomValidMeasurementUnit(ctx context.Context, db DBTX) (*GetRandomValidMeasurementUnitRow, error) {
	row := db.QueryRowContext(ctx, getRandomValidMeasurementUnit)
	var i GetRandomValidMeasurementUnitRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Volumetric,
		&i.IconPath,
		&i.Universal,
		&i.Metric,
		&i.Imperial,
		&i.Slug,
		&i.PluralName,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidMeasurementUnit = `-- name: GetValidMeasurementUnit :one
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
AND valid_measurement_units.id = $1
`

type GetValidMeasurementUnitRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidMeasurementUnit(ctx context.Context, db DBTX, id string) (*GetValidMeasurementUnitRow, error) {
	row := db.QueryRowContext(ctx, getValidMeasurementUnit, id)
	var i GetValidMeasurementUnitRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Volumetric,
		&i.IconPath,
		&i.Universal,
		&i.Metric,
		&i.Imperial,
		&i.Slug,
		&i.PluralName,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidMeasurementUnits = `-- name: GetValidMeasurementUnits :many
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at,
	(
		SELECT COUNT(valid_measurement_units.id)
		FROM valid_measurement_units
		WHERE valid_measurement_units.archived_at IS NULL
			AND valid_measurement_units.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_measurement_units.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_measurement_units.last_updated_at IS NULL
				OR valid_measurement_units.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_measurement_units.last_updated_at IS NULL
				OR valid_measurement_units.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(valid_measurement_units.id)
		FROM valid_measurement_units
		WHERE valid_measurement_units.archived_at IS NULL
	) AS total_count
FROM valid_measurement_units
WHERE
	valid_measurement_units.archived_at IS NULL
	AND valid_measurement_units.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_measurement_units.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_measurement_units.last_updated_at IS NULL
		OR valid_measurement_units.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_measurement_units.last_updated_at IS NULL
		OR valid_measurement_units.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
GROUP BY valid_measurement_units.id
ORDER BY valid_measurement_units.id
LIMIT $6
OFFSET $5
`

type GetValidMeasurementUnitsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetValidMeasurementUnitsRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
	FilteredCount int64
	TotalCount    int64
}

func (q *Queries) GetValidMeasurementUnits(ctx context.Context, db DBTX, arg *GetValidMeasurementUnitsParams) ([]*GetValidMeasurementUnitsRow, error) {
	rows, err := db.QueryContext(ctx, getValidMeasurementUnits,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidMeasurementUnitsRow{}
	for rows.Next() {
		var i GetValidMeasurementUnitsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
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

const getValidMeasurementUnitsNeedingIndexing = `-- name: GetValidMeasurementUnitsNeedingIndexing :many
SELECT valid_measurement_units.id
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	AND (
	valid_measurement_units.last_indexed_at IS NULL
	OR valid_measurement_units.last_indexed_at < NOW() - '24 hours'::INTERVAL
)
`

func (q *Queries) GetValidMeasurementUnitsNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getValidMeasurementUnitsNeedingIndexing)
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

const getValidMeasurementUnitsWithIDs = `-- name: GetValidMeasurementUnitsWithIDs :many
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	AND valid_measurement_units.id = ANY($1::text[])
`

type GetValidMeasurementUnitsWithIDsRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidMeasurementUnitsWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidMeasurementUnitsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidMeasurementUnitsWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidMeasurementUnitsWithIDsRow{}
	for rows.Next() {
		var i GetValidMeasurementUnitsWithIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
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

const searchForValidMeasurementUnits = `-- name: SearchForValidMeasurementUnits :many
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.name ILIKE '%' || $1::text || '%'
	AND valid_measurement_units.archived_at IS NULL
LIMIT 50
`

type SearchForValidMeasurementUnitsRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) SearchForValidMeasurementUnits(ctx context.Context, db DBTX, nameQuery string) ([]*SearchForValidMeasurementUnitsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidMeasurementUnits, nameQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidMeasurementUnitsRow{}
	for rows.Next() {
		var i SearchForValidMeasurementUnitsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
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

const searchValidMeasurementUnitsByIngredientID = `-- name: SearchValidMeasurementUnitsByIngredientID :many
SELECT
	DISTINCT(valid_measurement_units.id),
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.last_indexed_at,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at,
	(
		SELECT COUNT(valid_measurement_units.id)
		FROM valid_measurement_units
		WHERE valid_measurement_units.archived_at IS NULL
			AND valid_measurement_units.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_measurement_units.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_measurement_units.last_updated_at IS NULL
				OR valid_measurement_units.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_measurement_units.last_updated_at IS NULL
				OR valid_measurement_units.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (
				valid_ingredient_measurement_units.valid_ingredient_id = $5
				OR valid_measurement_units.universal = true
			)
	) AS filtered_count,
	(
		SELECT COUNT(valid_measurement_units.id)
		FROM valid_measurement_units
		WHERE valid_measurement_units.archived_at IS NULL
	) AS total_count
FROM valid_measurement_units
	JOIN valid_ingredient_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE
	(
		valid_ingredient_measurement_units.valid_ingredient_id = $5
		OR valid_measurement_units.universal = TRUE
	)
	AND valid_measurement_units.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_measurement_units.archived_at IS NULL
	AND valid_measurement_units.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_measurement_units.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_measurement_units.last_updated_at IS NULL
		OR valid_measurement_units.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_measurement_units.last_updated_at IS NULL
		OR valid_measurement_units.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT $7
OFFSET $6
`

type SearchValidMeasurementUnitsByIngredientIDParams struct {
	CreatedAfter      sql.NullTime
	CreatedBefore     sql.NullTime
	UpdatedBefore     sql.NullTime
	UpdatedAfter      sql.NullTime
	ValidIngredientID string
	QueryOffset       sql.NullInt32
	QueryLimit        sql.NullInt32
}

type SearchValidMeasurementUnitsByIngredientIDRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	LastIndexedAt sql.NullTime
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
	FilteredCount int64
	TotalCount    int64
}

func (q *Queries) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, db DBTX, arg *SearchValidMeasurementUnitsByIngredientIDParams) ([]*SearchValidMeasurementUnitsByIngredientIDRow, error) {
	rows, err := db.QueryContext(ctx, searchValidMeasurementUnitsByIngredientID,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.ValidIngredientID,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchValidMeasurementUnitsByIngredientIDRow{}
	for rows.Next() {
		var i SearchValidMeasurementUnitsByIngredientIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
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

const updateValidMeasurementUnit = `-- name: UpdateValidMeasurementUnit :execrows
UPDATE valid_measurement_units SET
	name = $1,
	description = $2,
	volumetric = $3,
	icon_path = $4,
	universal = $5,
	metric = $6,
	imperial = $7,
	slug = $8,
	plural_name = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $10
`

type UpdateValidMeasurementUnitParams struct {
	Name        string
	Description string
	Volumetric  sql.NullBool
	IconPath    string
	Universal   bool
	Metric      bool
	Imperial    bool
	Slug        string
	PluralName  string
	ID          string
}

func (q *Queries) UpdateValidMeasurementUnit(ctx context.Context, db DBTX, arg *UpdateValidMeasurementUnitParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidMeasurementUnit,
		arg.Name,
		arg.Description,
		arg.Volumetric,
		arg.IconPath,
		arg.Universal,
		arg.Metric,
		arg.Imperial,
		arg.Slug,
		arg.PluralName,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateValidMeasurementUnitLastIndexedAt = `-- name: UpdateValidMeasurementUnitLastIndexedAt :execrows
UPDATE valid_measurement_units SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateValidMeasurementUnitLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidMeasurementUnitLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
