// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_step_products_get_for_recipe.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetRecipeStepProductsForRecipe = `-- name: GetRecipeStepProductsForRecipe :many
SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
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
	valid_measurement_units.archived_at,
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.storage_instructions,
	recipe_step_products.created_at,
	recipe_step_products.last_updated_at,
	recipe_step_products.archived_at,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
	JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	LEFT OUTER JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1
`

type GetRecipeStepProductsForRecipeRow struct {
	ID                                 string
	Name                               string
	Type                               RecipeStepProductType
	ID_2                               sql.NullString
	Name_2                             sql.NullString
	Description                        sql.NullString
	Volumetric                         sql.NullBool
	IconPath                           sql.NullString
	Universal                          sql.NullBool
	Metric                             sql.NullBool
	Imperial                           sql.NullBool
	PluralName                         sql.NullString
	CreatedAt                          sql.NullTime
	LastUpdatedAt                      sql.NullTime
	ArchivedAt                         sql.NullTime
	MinimumQuantityValue               string
	MaximumQuantityValue               string
	QuantityNotes                      string
	Compostable                        bool
	MaximumStorageDurationInSeconds    sql.NullInt32
	MinimumStorageTemperatureInCelsius sql.NullString
	MaximumStorageTemperatureInCelsius sql.NullString
	StorageInstructions                string
	CreatedAt_2                        time.Time
	LastUpdatedAt_2                    sql.NullTime
	ArchivedAt_2                       sql.NullTime
	BelongsToRecipeStep                string
}

func (q *Queries) GetRecipeStepProductsForRecipe(ctx context.Context, belongsToRecipe string) ([]*GetRecipeStepProductsForRecipeRow, error) {
	rows, err := q.db.QueryContext(ctx, GetRecipeStepProductsForRecipe, belongsToRecipe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepProductsForRecipeRow{}
	for rows.Next() {
		var i GetRecipeStepProductsForRecipeRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.ID_2,
			&i.Name_2,
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
			&i.MinimumQuantityValue,
			&i.MaximumQuantityValue,
			&i.QuantityNotes,
			&i.Compostable,
			&i.MaximumStorageDurationInSeconds,
			&i.MinimumStorageTemperatureInCelsius,
			&i.MaximumStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
			&i.BelongsToRecipeStep,
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
