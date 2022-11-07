// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_ingredient_measurement_units_create.sql

package generated

import (
	"context"
)

const CreateValidIngredientMeasurementUnit = `-- name: CreateValidIngredientMeasurementUnit :exec
INSERT INTO valid_ingredient_measurement_units
(id,notes,valid_measurement_unit_id,valid_ingredient_id,minimum_allowable_quantity,maximum_allowable_quantity)
VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateValidIngredientMeasurementUnitParams struct {
	ID                       string `db:"id"`
	Notes                    string `db:"notes"`
	ValidMeasurementUnitID   string `db:"valid_measurement_unit_id"`
	ValidIngredientID        string `db:"valid_ingredient_id"`
	MinimumAllowableQuantity string `db:"minimum_allowable_quantity"`
	MaximumAllowableQuantity string `db:"maximum_allowable_quantity"`
}

func (q *Queries) CreateValidIngredientMeasurementUnit(ctx context.Context, db DBTX, arg *CreateValidIngredientMeasurementUnitParams) error {
	_, err := db.ExecContext(ctx, CreateValidIngredientMeasurementUnit,
		arg.ID,
		arg.Notes,
		arg.ValidMeasurementUnitID,
		arg.ValidIngredientID,
		arg.MinimumAllowableQuantity,
		arg.MaximumAllowableQuantity,
	)
	return err
}