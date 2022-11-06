// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_measurement_units_get_one.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetValidMeasurementUnit = `-- name: GetValidMeasurementUnit :exec
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
	PluralName    string
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidMeasurementUnit(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, GetValidMeasurementUnit, id)
	return err
}
