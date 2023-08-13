// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: list_all_by_recipe.sql

package generated

import (
	"context"
	"database/sql"
)

const listAllRecipePrepTasksByRecipe = `-- name: ListAllRecipePrepTasksByRecipe :many

SELECT
	recipe_prep_tasks.id,
	recipe_prep_tasks.name,
	recipe_prep_tasks.description,
	recipe_prep_tasks.notes,
	recipe_prep_tasks.optional,
	recipe_prep_tasks.explicit_storage_instructions,
	recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.storage_type,
	recipe_prep_tasks.minimum_storage_temperature_in_celsius,
	recipe_prep_tasks.maximum_storage_temperature_in_celsius,
	recipe_prep_tasks.belongs_to_recipe,
	recipe_prep_tasks.created_at,
	recipe_prep_tasks.last_updated_at,
	recipe_prep_tasks.archived_at,
	recipe_prep_task_steps.id,
	recipe_prep_task_steps.belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step
FROM recipe_prep_tasks
	 FULL OUTER JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	 FULL OUTER JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
	 FULL OUTER JOIN recipes ON recipe_prep_tasks.belongs_to_recipe=recipes.id
WHERE recipe_prep_tasks.archived_at IS NULL
  AND recipe_steps.archived_at IS NULL
  AND recipes.archived_at IS NULL
  AND recipes.id = $1
  AND recipe_steps.belongs_to_recipe = $1
`

type ListAllRecipePrepTasksByRecipeRow struct {
	CreatedAt                              sql.NullTime
	ArchivedAt                             sql.NullTime
	LastUpdatedAt                          sql.NullTime
	StorageType                            NullStorageContainerType
	Description                            sql.NullString
	ExplicitStorageInstructions            sql.NullString
	BelongsToRecipePrepTask                sql.NullString
	BelongsToRecipeStep                    sql.NullString
	ID                                     sql.NullString
	MinimumStorageTemperatureInCelsius     sql.NullString
	MaximumStorageTemperatureInCelsius     sql.NullString
	BelongsToRecipe                        sql.NullString
	Notes                                  sql.NullString
	ID_2                                   sql.NullString
	Name                                   sql.NullString
	MaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	MinimumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	Optional                               sql.NullBool
	SatisfiesRecipeStep                    sql.NullBool
}

func (q *Queries) ListAllRecipePrepTasksByRecipe(ctx context.Context, db DBTX, id string) ([]*ListAllRecipePrepTasksByRecipeRow, error) {
	rows, err := db.QueryContext(ctx, listAllRecipePrepTasksByRecipe, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListAllRecipePrepTasksByRecipeRow{}
	for rows.Next() {
		var i ListAllRecipePrepTasksByRecipeRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Notes,
			&i.Optional,
			&i.ExplicitStorageInstructions,
			&i.MinimumTimeBufferBeforeRecipeInSeconds,
			&i.MaximumTimeBufferBeforeRecipeInSeconds,
			&i.StorageType,
			&i.MinimumStorageTemperatureInCelsius,
			&i.MaximumStorageTemperatureInCelsius,
			&i.BelongsToRecipe,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.ID_2,
			&i.BelongsToRecipeStep,
			&i.BelongsToRecipePrepTask,
			&i.SatisfiesRecipeStep,
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
