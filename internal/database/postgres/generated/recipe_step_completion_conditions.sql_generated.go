// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: recipe_step_completion_conditions.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipeStepCompletionCondition = `-- name: ArchiveRecipeStepCompletionCondition :execrows

UPDATE recipe_step_completion_conditions SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2
`

type ArchiveRecipeStepCompletionConditionParams struct {
	BelongsToRecipeStep string
	ID                  string
}

func (q *Queries) ArchiveRecipeStepCompletionCondition(ctx context.Context, db DBTX, arg *ArchiveRecipeStepCompletionConditionParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveRecipeStepCompletionCondition, arg.BelongsToRecipeStep, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkRecipeStepCompletionConditionExistence = `-- name: CheckRecipeStepCompletionConditionExistence :one

SELECT EXISTS (
	SELECT recipe_step_completion_conditions.id
	FROM recipe_step_completion_conditions
		JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step=recipe_steps.id
		JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	WHERE recipe_step_completion_conditions.archived_at IS NULL
		AND recipe_step_completion_conditions.belongs_to_recipe_step = $1
		AND recipe_step_completion_conditions.id = $2
		AND recipe_steps.archived_at IS NULL
		AND recipe_steps.belongs_to_recipe = $3
		AND recipe_steps.id = $1
		AND recipes.archived_at IS NULL
		AND recipes.id = $3
)
`

type CheckRecipeStepCompletionConditionExistenceParams struct {
	RecipeStepID                    string
	RecipeStepCompletionConditionID string
	RecipeID                        string
}

func (q *Queries) CheckRecipeStepCompletionConditionExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepCompletionConditionExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepCompletionConditionExistence, arg.RecipeStepID, arg.RecipeStepCompletionConditionID, arg.RecipeID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipeStepCompletionCondition = `-- name: CreateRecipeStepCompletionCondition :exec

INSERT INTO recipe_step_completion_conditions (
	id,
	optional,
	notes,
	belongs_to_recipe_step,
	ingredient_state
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
`

type CreateRecipeStepCompletionConditionParams struct {
	ID                  string
	Notes               string
	BelongsToRecipeStep string
	IngredientState     string
	Optional            bool
}

func (q *Queries) CreateRecipeStepCompletionCondition(ctx context.Context, db DBTX, arg *CreateRecipeStepCompletionConditionParams) error {
	_, err := db.ExecContext(ctx, createRecipeStepCompletionCondition,
		arg.ID,
		arg.Optional,
		arg.Notes,
		arg.BelongsToRecipeStep,
		arg.IngredientState,
	)
	return err
}

const getAllRecipeStepCompletionConditionsForRecipe = `-- name: GetAllRecipeStepCompletionConditionsForRecipe :many

SELECT
	recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.optional,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.last_indexed_at as valid_ingredient_state_last_indexed_at,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at,
	recipe_step_completion_conditions.notes,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	recipe_step_completion_conditions.ingredient_state,
	recipe_step_completion_conditions.created_at,
	recipe_step_completion_conditions.last_updated_at,
	recipe_step_completion_conditions.archived_at,
	recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
	recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
	recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipes.archived_at IS NULL
	AND valid_ingredient_states.archived_at IS NULL
	AND recipes.id = $1
GROUP BY
	recipe_step_completion_conditions.id,
	recipe_step_completion_condition_ingredients.id,
	valid_ingredient_states.id
`

type GetAllRecipeStepCompletionConditionsForRecipeRow struct {
	RecipeStepCompletionConditionIngredientCreatedAt         time.Time
	CreatedAt                                                time.Time
	ValidIngredientStateCreatedAt                            time.Time
	ValidIngredientStateLastIndexedAt                        sql.NullTime
	RecipeStepCompletionConditionIngredientArchivedAt        sql.NullTime
	RecipeStepCompletionConditionIngredientLastUpdatedAt     sql.NullTime
	ArchivedAt                                               sql.NullTime
	LastUpdatedAt                                            sql.NullTime
	ValidIngredientStateArchivedAt                           sql.NullTime
	ValidIngredientStateLastUpdatedAt                        sql.NullTime
	ValidIngredientStateName                                 string
	BelongsToRecipeStep                                      string
	ValidIngredientStateIconPath                             string
	ValidIngredientStateDescription                          string
	ValidIngredientStateSlug                                 string
	ValidIngredientStatePastTense                            string
	Notes                                                    string
	ValidIngredientStateAttributeType                        IngredientAttributeType
	IngredientState                                          string
	RecipeStepCompletionConditionIngredientID                string
	ValidIngredientStateID                                   string
	RecipeStepCompletionConditionIngredientBelongsToRecipeS  string
	ID                                                       string
	RecipeStepCompletionConditionIngredientRecipeStepIngredi string
	Optional                                                 bool
}

func (q *Queries) GetAllRecipeStepCompletionConditionsForRecipe(ctx context.Context, db DBTX, id string) ([]*GetAllRecipeStepCompletionConditionsForRecipeRow, error) {
	rows, err := db.QueryContext(ctx, getAllRecipeStepCompletionConditionsForRecipe, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllRecipeStepCompletionConditionsForRecipeRow{}
	for rows.Next() {
		var i GetAllRecipeStepCompletionConditionsForRecipeRow
		if err := rows.Scan(
			&i.RecipeStepCompletionConditionIngredientID,
			&i.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
			&i.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			&i.ID,
			&i.Optional,
			&i.ValidIngredientStateID,
			&i.ValidIngredientStateName,
			&i.ValidIngredientStatePastTense,
			&i.ValidIngredientStateSlug,
			&i.ValidIngredientStateDescription,
			&i.ValidIngredientStateIconPath,
			&i.ValidIngredientStateAttributeType,
			&i.ValidIngredientStateLastIndexedAt,
			&i.ValidIngredientStateCreatedAt,
			&i.ValidIngredientStateLastUpdatedAt,
			&i.ValidIngredientStateArchivedAt,
			&i.Notes,
			&i.BelongsToRecipeStep,
			&i.IngredientState,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.RecipeStepCompletionConditionIngredientCreatedAt,
			&i.RecipeStepCompletionConditionIngredientLastUpdatedAt,
			&i.RecipeStepCompletionConditionIngredientArchivedAt,
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

const getRecipeStepCompletionConditionWithIngredients = `-- name: GetRecipeStepCompletionConditionWithIngredients :many

SELECT
	recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.optional,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.last_indexed_at as valid_ingredient_state_last_indexed_at,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at,
	recipe_step_completion_conditions.notes,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	recipe_step_completion_conditions.ingredient_state,
	recipe_step_completion_conditions.created_at,
	recipe_step_completion_conditions.last_updated_at,
	recipe_step_completion_conditions.archived_at,
	recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
	recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
	recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.archived_at IS NULL
	AND recipe_step_completion_conditions.belongs_to_recipe_step = $1
	AND recipe_step_completion_conditions.id = $2
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $3
	AND recipe_steps.id = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $3
`

type GetRecipeStepCompletionConditionWithIngredientsParams struct {
	RecipeStepID                    string
	RecipeStepCompletionConditionID string
	RecipeID                        string
}

type GetRecipeStepCompletionConditionWithIngredientsRow struct {
	RecipeStepCompletionConditionIngredientCreatedAt         time.Time
	CreatedAt                                                time.Time
	ValidIngredientStateCreatedAt                            time.Time
	ValidIngredientStateLastIndexedAt                        sql.NullTime
	RecipeStepCompletionConditionIngredientArchivedAt        sql.NullTime
	RecipeStepCompletionConditionIngredientLastUpdatedAt     sql.NullTime
	ArchivedAt                                               sql.NullTime
	LastUpdatedAt                                            sql.NullTime
	ValidIngredientStateArchivedAt                           sql.NullTime
	ValidIngredientStateLastUpdatedAt                        sql.NullTime
	ValidIngredientStateName                                 string
	BelongsToRecipeStep                                      string
	ValidIngredientStateIconPath                             string
	ValidIngredientStateDescription                          string
	ValidIngredientStateSlug                                 string
	ValidIngredientStatePastTense                            string
	Notes                                                    string
	ValidIngredientStateAttributeType                        IngredientAttributeType
	IngredientState                                          string
	RecipeStepCompletionConditionIngredientID                string
	ValidIngredientStateID                                   string
	RecipeStepCompletionConditionIngredientBelongsToRecipeS  string
	ID                                                       string
	RecipeStepCompletionConditionIngredientRecipeStepIngredi string
	Optional                                                 bool
}

func (q *Queries) GetRecipeStepCompletionConditionWithIngredients(ctx context.Context, db DBTX, arg *GetRecipeStepCompletionConditionWithIngredientsParams) ([]*GetRecipeStepCompletionConditionWithIngredientsRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeStepCompletionConditionWithIngredients, arg.RecipeStepID, arg.RecipeStepCompletionConditionID, arg.RecipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepCompletionConditionWithIngredientsRow{}
	for rows.Next() {
		var i GetRecipeStepCompletionConditionWithIngredientsRow
		if err := rows.Scan(
			&i.RecipeStepCompletionConditionIngredientID,
			&i.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
			&i.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			&i.ID,
			&i.Optional,
			&i.ValidIngredientStateID,
			&i.ValidIngredientStateName,
			&i.ValidIngredientStatePastTense,
			&i.ValidIngredientStateSlug,
			&i.ValidIngredientStateDescription,
			&i.ValidIngredientStateIconPath,
			&i.ValidIngredientStateAttributeType,
			&i.ValidIngredientStateLastIndexedAt,
			&i.ValidIngredientStateCreatedAt,
			&i.ValidIngredientStateLastUpdatedAt,
			&i.ValidIngredientStateArchivedAt,
			&i.Notes,
			&i.BelongsToRecipeStep,
			&i.IngredientState,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.RecipeStepCompletionConditionIngredientCreatedAt,
			&i.RecipeStepCompletionConditionIngredientLastUpdatedAt,
			&i.RecipeStepCompletionConditionIngredientArchivedAt,
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

const getRecipeStepCompletionConditions = `-- name: GetRecipeStepCompletionConditions :many

SELECT
	recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.optional,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.last_indexed_at as valid_ingredient_state_last_indexed_at,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at,
	recipe_step_completion_conditions.notes,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	recipe_step_completion_conditions.ingredient_state,
	recipe_step_completion_conditions.created_at,
	recipe_step_completion_conditions.last_updated_at,
	recipe_step_completion_conditions.archived_at,
	recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
	recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
	recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at,
	(
		SELECT COUNT(recipe_step_completion_condition_ingredients.id)
		FROM recipe_step_completion_condition_ingredients
		WHERE recipe_step_completion_condition_ingredients.archived_at IS NULL
			AND recipe_step_completion_condition_ingredients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND recipe_step_completion_condition_ingredients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				recipe_step_completion_condition_ingredients.last_updated_at IS NULL
				OR recipe_step_completion_condition_ingredients.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				recipe_step_completion_condition_ingredients.last_updated_at IS NULL
				OR recipe_step_completion_condition_ingredients.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND recipe_step_completion_conditions.belongs_to_recipe_step = $5
	) AS filtered_count,
	(
		SELECT COUNT(recipe_step_completion_condition_ingredients.id)
		FROM recipe_step_completion_condition_ingredients
		WHERE recipe_step_completion_condition_ingredients.archived_at IS NULL
	) AS total_count
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND recipe_step_completion_condition_ingredients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		recipe_step_completion_condition_ingredients.last_updated_at IS NULL
		OR recipe_step_completion_condition_ingredients.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		recipe_step_completion_condition_ingredients.last_updated_at IS NULL
		OR recipe_step_completion_condition_ingredients.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
	AND recipe_step_completion_conditions.belongs_to_recipe_step = $5
LIMIT $7
OFFSET $6
`

type GetRecipeStepCompletionConditionsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	RecipeStepID  string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipeStepCompletionConditionsRow struct {
	RecipeStepCompletionConditionIngredientCreatedAt         time.Time
	CreatedAt                                                time.Time
	ValidIngredientStateCreatedAt                            time.Time
	ValidIngredientStateLastIndexedAt                        sql.NullTime
	RecipeStepCompletionConditionIngredientArchivedAt        sql.NullTime
	RecipeStepCompletionConditionIngredientLastUpdatedAt     sql.NullTime
	ArchivedAt                                               sql.NullTime
	LastUpdatedAt                                            sql.NullTime
	ValidIngredientStateArchivedAt                           sql.NullTime
	ValidIngredientStateLastUpdatedAt                        sql.NullTime
	ValidIngredientStateName                                 string
	RecipeStepCompletionConditionIngredientID                string
	ValidIngredientStateIconPath                             string
	ValidIngredientStateDescription                          string
	ValidIngredientStateSlug                                 string
	ValidIngredientStatePastTense                            string
	Notes                                                    string
	BelongsToRecipeStep                                      string
	IngredientState                                          string
	ValidIngredientStateAttributeType                        IngredientAttributeType
	ValidIngredientStateID                                   string
	RecipeStepCompletionConditionIngredientBelongsToRecipeS  string
	ID                                                       string
	RecipeStepCompletionConditionIngredientRecipeStepIngredi string
	FilteredCount                                            int64
	TotalCount                                               int64
	Optional                                                 bool
}

func (q *Queries) GetRecipeStepCompletionConditions(ctx context.Context, db DBTX, arg *GetRecipeStepCompletionConditionsParams) ([]*GetRecipeStepCompletionConditionsRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeStepCompletionConditions,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.RecipeStepID,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepCompletionConditionsRow{}
	for rows.Next() {
		var i GetRecipeStepCompletionConditionsRow
		if err := rows.Scan(
			&i.RecipeStepCompletionConditionIngredientID,
			&i.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
			&i.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			&i.ID,
			&i.Optional,
			&i.ValidIngredientStateID,
			&i.ValidIngredientStateName,
			&i.ValidIngredientStatePastTense,
			&i.ValidIngredientStateSlug,
			&i.ValidIngredientStateDescription,
			&i.ValidIngredientStateIconPath,
			&i.ValidIngredientStateAttributeType,
			&i.ValidIngredientStateLastIndexedAt,
			&i.ValidIngredientStateCreatedAt,
			&i.ValidIngredientStateLastUpdatedAt,
			&i.ValidIngredientStateArchivedAt,
			&i.Notes,
			&i.BelongsToRecipeStep,
			&i.IngredientState,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.RecipeStepCompletionConditionIngredientCreatedAt,
			&i.RecipeStepCompletionConditionIngredientLastUpdatedAt,
			&i.RecipeStepCompletionConditionIngredientArchivedAt,
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

const updateRecipeStepCompletionCondition = `-- name: UpdateRecipeStepCompletionCondition :execrows

UPDATE recipe_step_completion_conditions SET
	optional = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	ingredient_state = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $5
`

type UpdateRecipeStepCompletionConditionParams struct {
	Notes               string
	BelongsToRecipeStep string
	IngredientState     string
	ID                  string
	Optional            bool
}

func (q *Queries) UpdateRecipeStepCompletionCondition(ctx context.Context, db DBTX, arg *UpdateRecipeStepCompletionConditionParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipeStepCompletionCondition,
		arg.Optional,
		arg.Notes,
		arg.BelongsToRecipeStep,
		arg.IngredientState,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
