// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_step_instruments_update.sql

package generated

import (
	"context"
	"database/sql"
)

const UpdateRecipeStepInstrument = `-- name: UpdateRecipeStepInstrument :exec
UPDATE recipe_step_instruments SET
	instrument_id = $1,
	recipe_step_product_id = $2,
	name = $3,
	product_of_recipe_step = $4,
	notes = $5,
	preference_rank = $6,
	optional = $7,
	minimum_quantity = $8,
	maximum_quantity = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $10
	AND id = $11
`

type UpdateRecipeStepInstrumentParams struct {
	Name                string         `db:"name"`
	BelongsToRecipeStep string         `db:"belongs_to_recipe_step"`
	Notes               string         `db:"notes"`
	ID                  string         `db:"id"`
	InstrumentID        sql.NullString `db:"instrument_id"`
	RecipeStepProductID sql.NullString `db:"recipe_step_product_id"`
	PreferenceRank      int32          `db:"preference_rank"`
	MinimumQuantity     int32          `db:"minimum_quantity"`
	MaximumQuantity     int32          `db:"maximum_quantity"`
	Optional            bool           `db:"optional"`
	ProductOfRecipeStep bool           `db:"product_of_recipe_step"`
}

func (q *Queries) UpdateRecipeStepInstrument(ctx context.Context, db DBTX, arg *UpdateRecipeStepInstrumentParams) error {
	_, err := db.ExecContext(ctx, UpdateRecipeStepInstrument,
		arg.InstrumentID,
		arg.RecipeStepProductID,
		arg.Name,
		arg.ProductOfRecipeStep,
		arg.Notes,
		arg.PreferenceRank,
		arg.Optional,
		arg.MinimumQuantity,
		arg.MaximumQuantity,
		arg.BelongsToRecipeStep,
		arg.ID,
	)
	return err
}