// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_measurement_units_get_random.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetRandomValidMeasurementUnit = `-- name: GetRandomValidMeasurementUnit :one
SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	ORDER BY random() LIMIT 1
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
	PluralName    string
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetRandomValidMeasurementUnit(ctx context.Context) (*GetRandomValidMeasurementUnitRow, error) {
	row := q.db.QueryRowContext(ctx, GetRandomValidMeasurementUnit)
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
		&i.PluralName,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}
