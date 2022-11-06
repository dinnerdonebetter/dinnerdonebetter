// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_measurement_units_exists.sql

package generated

import (
	"context"
)

const ValidMeasurementUnitExists = `-- name: ValidMeasurementUnitExists :exec
SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_at IS NULL AND valid_measurement_units.id = $1 )
`

func (q *Queries) ValidMeasurementUnitExists(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, ValidMeasurementUnitExists, id)
	return err
}
