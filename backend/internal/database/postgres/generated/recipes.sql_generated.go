// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: recipes.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipe = `-- name: ArchiveRecipe :execrows
UPDATE recipes SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2
`

type ArchiveRecipeParams struct {
	CreatedByUser string
	ID            string
}

func (q *Queries) ArchiveRecipe(ctx context.Context, db DBTX, arg *ArchiveRecipeParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveRecipe, arg.CreatedByUser, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkRecipeExistence = `-- name: CheckRecipeExistence :one
SELECT EXISTS (
	SELECT recipes.id
	FROM recipes
	WHERE recipes.archived_at IS NULL
		AND recipes.id = $1
)
`

func (q *Queries) CheckRecipeExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipe = `-- name: CreateRecipe :exec
INSERT INTO recipes (
	id,
	name,
	slug,
	source,
	description,
	inspired_by_recipe_id,
	min_estimated_portions,
	max_estimated_portions,
	portion_name,
	plural_portion_name,
	seal_of_approval,
	eligible_for_meals,
	yields_component_type,
	created_by_user
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
	$10,
	$11,
	$12,
	$13,
	$14
)
`

type CreateRecipeParams struct {
	MinEstimatedPortions string
	ID                   string
	Slug                 string
	Source               string
	Description          string
	CreatedByUser        string
	Name                 string
	YieldsComponentType  ComponentType
	PortionName          string
	PluralPortionName    string
	MaxEstimatedPortions sql.NullString
	InspiredByRecipeID   sql.NullString
	SealOfApproval       bool
	EligibleForMeals     bool
}

func (q *Queries) CreateRecipe(ctx context.Context, db DBTX, arg *CreateRecipeParams) error {
	_, err := db.ExecContext(ctx, createRecipe,
		arg.ID,
		arg.Name,
		arg.Slug,
		arg.Source,
		arg.Description,
		arg.InspiredByRecipeID,
		arg.MinEstimatedPortions,
		arg.MaxEstimatedPortions,
		arg.PortionName,
		arg.PluralPortionName,
		arg.SealOfApproval,
		arg.EligibleForMeals,
		arg.YieldsComponentType,
		arg.CreatedByUser,
	)
	return err
}

const getRecipeByID = `-- name: GetRecipeByID :many
SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.last_indexed_at,
	recipes.last_validated_at,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	recipe_steps.id as recipe_step_id,
	recipe_steps.index as recipe_step_index,
	valid_preparations.id as recipe_step_preparation_id,
	valid_preparations.name as recipe_step_preparation_name,
	valid_preparations.description as recipe_step_preparation_description,
	valid_preparations.icon_path as recipe_step_preparation_icon_path,
	valid_preparations.yields_nothing as recipe_step_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as recipe_step_preparation_restrict_to_ingredients,
	valid_preparations.past_tense as recipe_step_preparation_past_tense,
	valid_preparations.slug as recipe_step_preparation_slug,
	valid_preparations.minimum_ingredient_count as recipe_step_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as recipe_step_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as recipe_step_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as recipe_step_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as recipe_step_preparation_temperature_required,
	valid_preparations.time_estimate_required as recipe_step_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as recipe_step_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as recipe_step_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as recipe_step_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as recipe_step_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as recipe_step_preparation_maximum_vessel_count,
	valid_preparations.last_indexed_at as recipe_step_preparation_last_indexed_at,
	valid_preparations.created_at as recipe_step_preparation_created_at,
	valid_preparations.last_updated_at as recipe_step_preparation_last_updated_at,
	valid_preparations.archived_at as recipe_step_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds as recipe_step_minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds as recipe_step_maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius as recipe_step_minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius as recipe_step_maximum_temperature_in_celsius,
	recipe_steps.notes as recipe_step_notes,
	recipe_steps.explicit_instructions as recipe_step_explicit_instructions,
	recipe_steps.condition_expression as recipe_step_condition_expression,
	recipe_steps.optional as recipe_step_optional,
	recipe_steps.start_timer_automatically as recipe_step_start_timer_automatically,
	recipe_steps.created_at as recipe_step_created_at,
	recipe_steps.last_updated_at as recipe_step_last_updated_at,
	recipe_steps.archived_at as recipe_step_archived_at,
	recipe_steps.belongs_to_recipe as recipe_step_belongs_to_recipe
FROM recipes
	JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
	AND recipes.id = $1
ORDER BY recipe_steps.index
`

type GetRecipeByIDRow struct {
	CreatedAt                                        time.Time
	RecipeStepCreatedAt                              time.Time
	RecipeStepPreparationCreatedAt                   time.Time
	RecipeStepPreparationLastUpdatedAt               sql.NullTime
	LastIndexedAt                                    sql.NullTime
	RecipeStepArchivedAt                             sql.NullTime
	RecipeStepLastUpdatedAt                          sql.NullTime
	RecipeStepPreparationArchivedAt                  sql.NullTime
	RecipeStepPreparationLastIndexedAt               sql.NullTime
	ArchivedAt                                       sql.NullTime
	LastUpdatedAt                                    sql.NullTime
	LastValidatedAt                                  sql.NullTime
	ID                                               string
	PortionName                                      string
	RecipeStepBelongsToRecipe                        string
	Source                                           string
	MinEstimatedPortions                             string
	PluralPortionName                                string
	CreatedByUser                                    string
	RecipeStepID                                     string
	Name                                             string
	RecipeStepPreparationID                          string
	RecipeStepPreparationName                        string
	RecipeStepPreparationDescription                 string
	RecipeStepPreparationIconPath                    string
	RecipeStepConditionExpression                    string
	RecipeStepExplicitInstructions                   string
	RecipeStepPreparationPastTense                   string
	RecipeStepPreparationSlug                        string
	RecipeStepNotes                                  string
	YieldsComponentType                              ComponentType
	Description                                      string
	Slug                                             string
	RecipeStepMinimumTemperatureInCelsius            sql.NullString
	RecipeStepMaximumTemperatureInCelsius            sql.NullString
	InspiredByRecipeID                               sql.NullString
	MaxEstimatedPortions                             sql.NullString
	RecipeStepMaximumEstimatedTimeInSeconds          sql.NullInt64
	RecipeStepMinimumEstimatedTimeInSeconds          sql.NullInt64
	RecipeStepPreparationMaximumVesselCount          sql.NullInt32
	RecipeStepPreparationMaximumInstrumentCount      sql.NullInt32
	RecipeStepPreparationMaximumIngredientCount      sql.NullInt32
	RecipeStepPreparationMinimumInstrumentCount      int32
	RecipeStepPreparationMinimumIngredientCount      int32
	RecipeStepPreparationMinimumVesselCount          int32
	RecipeStepIndex                                  int32
	RecipeStepPreparationOnlyForVessels              bool
	RecipeStepStartTimerAutomatically                bool
	RecipeStepPreparationConditionExpressionRequired bool
	RecipeStepPreparationRestrictToIngredients       bool
	RecipeStepPreparationYieldsNothing               bool
	RecipeStepOptional                               bool
	RecipeStepPreparationConsumesVessel              bool
	RecipeStepPreparationTemperatureRequired         bool
	SealOfApproval                                   bool
	RecipeStepPreparationTimeEstimateRequired        bool
	EligibleForMeals                                 bool
}

func (q *Queries) GetRecipeByID(ctx context.Context, db DBTX, recipeID string) ([]*GetRecipeByIDRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeByID, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeByIDRow{}
	for rows.Next() {
		var i GetRecipeByIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Source,
			&i.Description,
			&i.InspiredByRecipeID,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.PortionName,
			&i.PluralPortionName,
			&i.SealOfApproval,
			&i.EligibleForMeals,
			&i.YieldsComponentType,
			&i.LastIndexedAt,
			&i.LastValidatedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
			&i.RecipeStepID,
			&i.RecipeStepIndex,
			&i.RecipeStepPreparationID,
			&i.RecipeStepPreparationName,
			&i.RecipeStepPreparationDescription,
			&i.RecipeStepPreparationIconPath,
			&i.RecipeStepPreparationYieldsNothing,
			&i.RecipeStepPreparationRestrictToIngredients,
			&i.RecipeStepPreparationPastTense,
			&i.RecipeStepPreparationSlug,
			&i.RecipeStepPreparationMinimumIngredientCount,
			&i.RecipeStepPreparationMaximumIngredientCount,
			&i.RecipeStepPreparationMinimumInstrumentCount,
			&i.RecipeStepPreparationMaximumInstrumentCount,
			&i.RecipeStepPreparationTemperatureRequired,
			&i.RecipeStepPreparationTimeEstimateRequired,
			&i.RecipeStepPreparationConditionExpressionRequired,
			&i.RecipeStepPreparationConsumesVessel,
			&i.RecipeStepPreparationOnlyForVessels,
			&i.RecipeStepPreparationMinimumVesselCount,
			&i.RecipeStepPreparationMaximumVesselCount,
			&i.RecipeStepPreparationLastIndexedAt,
			&i.RecipeStepPreparationCreatedAt,
			&i.RecipeStepPreparationLastUpdatedAt,
			&i.RecipeStepPreparationArchivedAt,
			&i.RecipeStepMinimumEstimatedTimeInSeconds,
			&i.RecipeStepMaximumEstimatedTimeInSeconds,
			&i.RecipeStepMinimumTemperatureInCelsius,
			&i.RecipeStepMaximumTemperatureInCelsius,
			&i.RecipeStepNotes,
			&i.RecipeStepExplicitInstructions,
			&i.RecipeStepConditionExpression,
			&i.RecipeStepOptional,
			&i.RecipeStepStartTimerAutomatically,
			&i.RecipeStepCreatedAt,
			&i.RecipeStepLastUpdatedAt,
			&i.RecipeStepArchivedAt,
			&i.RecipeStepBelongsToRecipe,
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

const getRecipeByIDAndAuthorID = `-- name: GetRecipeByIDAndAuthorID :many
SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.last_indexed_at,
	recipes.last_validated_at,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	recipe_steps.id as recipe_step_id,
	recipe_steps.index as recipe_step_index,
	valid_preparations.id as recipe_step_preparation_id,
	valid_preparations.name as recipe_step_preparation_name,
	valid_preparations.description as recipe_step_preparation_description,
	valid_preparations.icon_path as recipe_step_preparation_icon_path,
	valid_preparations.yields_nothing as recipe_step_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as recipe_step_preparation_restrict_to_ingredients,
	valid_preparations.past_tense as recipe_step_preparation_past_tense,
	valid_preparations.slug as recipe_step_preparation_slug,
	valid_preparations.minimum_ingredient_count as recipe_step_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as recipe_step_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as recipe_step_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as recipe_step_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as recipe_step_preparation_temperature_required,
	valid_preparations.time_estimate_required as recipe_step_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as recipe_step_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as recipe_step_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as recipe_step_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as recipe_step_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as recipe_step_preparation_maximum_vessel_count,
	valid_preparations.last_indexed_at as recipe_step_preparation_last_indexed_at,
	valid_preparations.created_at as recipe_step_preparation_created_at,
	valid_preparations.last_updated_at as recipe_step_preparation_last_updated_at,
	valid_preparations.archived_at as recipe_step_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds as recipe_step_minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds as recipe_step_maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius as recipe_step_minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius as recipe_step_maximum_temperature_in_celsius,
	recipe_steps.notes as recipe_step_notes,
	recipe_steps.explicit_instructions as recipe_step_explicit_instructions,
	recipe_steps.condition_expression as recipe_step_condition_expression,
	recipe_steps.optional as recipe_step_optional,
	recipe_steps.start_timer_automatically as recipe_step_start_timer_automatically,
	recipe_steps.created_at as recipe_step_created_at,
	recipe_steps.last_updated_at as recipe_step_last_updated_at,
	recipe_steps.archived_at as recipe_step_archived_at,
	recipe_steps.belongs_to_recipe as recipe_step_belongs_to_recipe
FROM recipes
	FULL OUTER JOIN recipe_steps ON recipes.id=recipe_steps.belongs_to_recipe
	FULL OUTER JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipes.archived_at IS NULL
	AND recipes.id = $1
	AND recipes.created_by_user = $2
ORDER BY recipe_steps.index
`

type GetRecipeByIDAndAuthorIDParams struct {
	RecipeID      string
	CreatedByUser string
}

type GetRecipeByIDAndAuthorIDRow struct {
	RecipeStepPreparationLastUpdatedAt               sql.NullTime
	LastIndexedAt                                    sql.NullTime
	RecipeStepArchivedAt                             sql.NullTime
	RecipeStepLastUpdatedAt                          sql.NullTime
	RecipeStepCreatedAt                              sql.NullTime
	RecipeStepPreparationArchivedAt                  sql.NullTime
	RecipeStepPreparationCreatedAt                   sql.NullTime
	RecipeStepPreparationLastIndexedAt               sql.NullTime
	ArchivedAt                                       sql.NullTime
	LastUpdatedAt                                    sql.NullTime
	CreatedAt                                        sql.NullTime
	LastValidatedAt                                  sql.NullTime
	MinEstimatedPortions                             sql.NullString
	RecipeStepMinimumTemperatureInCelsius            sql.NullString
	RecipeStepBelongsToRecipe                        sql.NullString
	Slug                                             sql.NullString
	PluralPortionName                                sql.NullString
	PortionName                                      sql.NullString
	CreatedByUser                                    sql.NullString
	Source                                           sql.NullString
	Description                                      sql.NullString
	RecipeStepPreparationID                          sql.NullString
	RecipeStepPreparationName                        sql.NullString
	RecipeStepPreparationDescription                 sql.NullString
	RecipeStepPreparationIconPath                    sql.NullString
	RecipeStepConditionExpression                    sql.NullString
	RecipeStepExplicitInstructions                   sql.NullString
	RecipeStepPreparationPastTense                   sql.NullString
	RecipeStepPreparationSlug                        sql.NullString
	Name                                             sql.NullString
	RecipeStepID                                     sql.NullString
	RecipeStepNotes                                  sql.NullString
	RecipeStepMaximumTemperatureInCelsius            sql.NullString
	YieldsComponentType                              NullComponentType
	InspiredByRecipeID                               sql.NullString
	ID                                               sql.NullString
	MaxEstimatedPortions                             sql.NullString
	RecipeStepMinimumEstimatedTimeInSeconds          sql.NullInt64
	RecipeStepMaximumEstimatedTimeInSeconds          sql.NullInt64
	RecipeStepPreparationMaximumInstrumentCount      sql.NullInt32
	RecipeStepPreparationMinimumVesselCount          sql.NullInt32
	RecipeStepPreparationMaximumVesselCount          sql.NullInt32
	RecipeStepPreparationMaximumIngredientCount      sql.NullInt32
	RecipeStepPreparationMinimumIngredientCount      sql.NullInt32
	RecipeStepIndex                                  sql.NullInt32
	RecipeStepPreparationMinimumInstrumentCount      sql.NullInt32
	RecipeStepPreparationConditionExpressionRequired sql.NullBool
	RecipeStepPreparationTemperatureRequired         sql.NullBool
	RecipeStepPreparationTimeEstimateRequired        sql.NullBool
	RecipeStepPreparationRestrictToIngredients       sql.NullBool
	RecipeStepPreparationYieldsNothing               sql.NullBool
	RecipeStepOptional                               sql.NullBool
	RecipeStepStartTimerAutomatically                sql.NullBool
	RecipeStepPreparationOnlyForVessels              sql.NullBool
	RecipeStepPreparationConsumesVessel              sql.NullBool
	SealOfApproval                                   sql.NullBool
	EligibleForMeals                                 sql.NullBool
}

func (q *Queries) GetRecipeByIDAndAuthorID(ctx context.Context, db DBTX, arg *GetRecipeByIDAndAuthorIDParams) ([]*GetRecipeByIDAndAuthorIDRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeByIDAndAuthorID, arg.RecipeID, arg.CreatedByUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeByIDAndAuthorIDRow{}
	for rows.Next() {
		var i GetRecipeByIDAndAuthorIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Source,
			&i.Description,
			&i.InspiredByRecipeID,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.PortionName,
			&i.PluralPortionName,
			&i.SealOfApproval,
			&i.EligibleForMeals,
			&i.YieldsComponentType,
			&i.LastIndexedAt,
			&i.LastValidatedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
			&i.RecipeStepID,
			&i.RecipeStepIndex,
			&i.RecipeStepPreparationID,
			&i.RecipeStepPreparationName,
			&i.RecipeStepPreparationDescription,
			&i.RecipeStepPreparationIconPath,
			&i.RecipeStepPreparationYieldsNothing,
			&i.RecipeStepPreparationRestrictToIngredients,
			&i.RecipeStepPreparationPastTense,
			&i.RecipeStepPreparationSlug,
			&i.RecipeStepPreparationMinimumIngredientCount,
			&i.RecipeStepPreparationMaximumIngredientCount,
			&i.RecipeStepPreparationMinimumInstrumentCount,
			&i.RecipeStepPreparationMaximumInstrumentCount,
			&i.RecipeStepPreparationTemperatureRequired,
			&i.RecipeStepPreparationTimeEstimateRequired,
			&i.RecipeStepPreparationConditionExpressionRequired,
			&i.RecipeStepPreparationConsumesVessel,
			&i.RecipeStepPreparationOnlyForVessels,
			&i.RecipeStepPreparationMinimumVesselCount,
			&i.RecipeStepPreparationMaximumVesselCount,
			&i.RecipeStepPreparationLastIndexedAt,
			&i.RecipeStepPreparationCreatedAt,
			&i.RecipeStepPreparationLastUpdatedAt,
			&i.RecipeStepPreparationArchivedAt,
			&i.RecipeStepMinimumEstimatedTimeInSeconds,
			&i.RecipeStepMaximumEstimatedTimeInSeconds,
			&i.RecipeStepMinimumTemperatureInCelsius,
			&i.RecipeStepMaximumTemperatureInCelsius,
			&i.RecipeStepNotes,
			&i.RecipeStepExplicitInstructions,
			&i.RecipeStepConditionExpression,
			&i.RecipeStepOptional,
			&i.RecipeStepStartTimerAutomatically,
			&i.RecipeStepCreatedAt,
			&i.RecipeStepLastUpdatedAt,
			&i.RecipeStepArchivedAt,
			&i.RecipeStepBelongsToRecipe,
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

const getRecipeIDsForMeal = `-- name: GetRecipeIDsForMeal :many
SELECT recipes.id
FROM recipes
	JOIN meal_components ON meal_components.recipe_id = recipes.id
	JOIN meals ON meal_components.meal_id = meals.id
WHERE
	recipes.archived_at IS NULL
	AND meals.id = $1
GROUP BY recipes.id
ORDER BY recipes.id
`

func (q *Queries) GetRecipeIDsForMeal(ctx context.Context, db DBTX, mealID string) ([]string, error) {
	rows, err := db.QueryContext(ctx, getRecipeIDsForMeal, mealID)
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

const getRecipes = `-- name: GetRecipes :many
SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.last_indexed_at,
	recipes.last_validated_at,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
			AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
	) AS total_count
FROM recipes
	WHERE recipes.archived_at IS NULL
	AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT $6
OFFSET $5
`

type GetRecipesParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipesRow struct {
	CreatedAt            time.Time
	LastValidatedAt      sql.NullTime
	LastIndexedAt        sql.NullTime
	LastUpdatedAt        sql.NullTime
	ArchivedAt           sql.NullTime
	MinEstimatedPortions string
	ID                   string
	CreatedByUser        string
	PortionName          string
	PluralPortionName    string
	Description          string
	Source               string
	YieldsComponentType  ComponentType
	Slug                 string
	Name                 string
	InspiredByRecipeID   sql.NullString
	MaxEstimatedPortions sql.NullString
	FilteredCount        int64
	TotalCount           int64
	EligibleForMeals     bool
	SealOfApproval       bool
}

func (q *Queries) GetRecipes(ctx context.Context, db DBTX, arg *GetRecipesParams) ([]*GetRecipesRow, error) {
	rows, err := db.QueryContext(ctx, getRecipes,
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
	items := []*GetRecipesRow{}
	for rows.Next() {
		var i GetRecipesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Source,
			&i.Description,
			&i.InspiredByRecipeID,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.PortionName,
			&i.PluralPortionName,
			&i.SealOfApproval,
			&i.EligibleForMeals,
			&i.YieldsComponentType,
			&i.LastIndexedAt,
			&i.LastValidatedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
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

const getRecipesCreatedByUser = `-- name: GetRecipesCreatedByUser :many
SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.last_indexed_at,
	recipes.last_validated_at,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
			AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND recipes.created_by_user = $5
	) AS filtered_count,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
			AND recipes.created_by_user = $5
	) AS total_count
FROM recipes
	WHERE recipes.archived_at IS NULL AND
	recipes.created_by_user = $5
	AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
	AND recipes.created_by_user = $5
LIMIT $7
OFFSET $6
`

type GetRecipesCreatedByUserParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	CreatedByUser string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipesCreatedByUserRow struct {
	CreatedAt            time.Time
	LastValidatedAt      sql.NullTime
	LastIndexedAt        sql.NullTime
	LastUpdatedAt        sql.NullTime
	ArchivedAt           sql.NullTime
	MinEstimatedPortions string
	ID                   string
	CreatedByUser        string
	PortionName          string
	PluralPortionName    string
	Description          string
	Source               string
	YieldsComponentType  ComponentType
	Slug                 string
	Name                 string
	InspiredByRecipeID   sql.NullString
	MaxEstimatedPortions sql.NullString
	FilteredCount        int64
	TotalCount           int64
	EligibleForMeals     bool
	SealOfApproval       bool
}

func (q *Queries) GetRecipesCreatedByUser(ctx context.Context, db DBTX, arg *GetRecipesCreatedByUserParams) ([]*GetRecipesCreatedByUserRow, error) {
	rows, err := db.QueryContext(ctx, getRecipesCreatedByUser,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.CreatedByUser,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipesCreatedByUserRow{}
	for rows.Next() {
		var i GetRecipesCreatedByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Source,
			&i.Description,
			&i.InspiredByRecipeID,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.PortionName,
			&i.PluralPortionName,
			&i.SealOfApproval,
			&i.EligibleForMeals,
			&i.YieldsComponentType,
			&i.LastIndexedAt,
			&i.LastValidatedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
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

const getRecipesNeedingIndexing = `-- name: GetRecipesNeedingIndexing :many
SELECT recipes.id
FROM recipes
WHERE recipes.archived_at IS NULL
	AND (
		recipes.last_indexed_at IS NULL
		OR recipes.last_indexed_at < NOW() - '24 hours'::INTERVAL
	)
`

func (q *Queries) GetRecipesNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getRecipesNeedingIndexing)
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

const recipeSearch = `-- name: RecipeSearch :many
SELECT
	recipes.id,
	recipes.name,
	recipes.slug,
	recipes.source,
	recipes.description,
	recipes.inspired_by_recipe_id,
	recipes.min_estimated_portions,
	recipes.max_estimated_portions,
	recipes.portion_name,
	recipes.plural_portion_name,
	recipes.seal_of_approval,
	recipes.eligible_for_meals,
	recipes.yields_component_type,
	recipes.last_indexed_at,
	recipes.last_validated_at,
	recipes.created_at,
	recipes.last_updated_at,
	recipes.archived_at,
	recipes.created_by_user,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
			AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				recipes.last_updated_at IS NULL
				OR recipes.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(recipes.id)
		FROM recipes
		WHERE recipes.archived_at IS NULL
	) AS total_count
FROM recipes
WHERE recipes.archived_at IS NULL
	AND recipes.name ILIKE '%' || $5::text || '%'
	AND recipes.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND recipes.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		recipes.last_updated_at IS NULL
		OR recipes.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT $7
OFFSET $6
`

type RecipeSearchParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	Query         string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type RecipeSearchRow struct {
	CreatedAt            time.Time
	LastValidatedAt      sql.NullTime
	LastIndexedAt        sql.NullTime
	LastUpdatedAt        sql.NullTime
	ArchivedAt           sql.NullTime
	MinEstimatedPortions string
	ID                   string
	CreatedByUser        string
	PortionName          string
	PluralPortionName    string
	Description          string
	Source               string
	YieldsComponentType  ComponentType
	Slug                 string
	Name                 string
	InspiredByRecipeID   sql.NullString
	MaxEstimatedPortions sql.NullString
	FilteredCount        int64
	TotalCount           int64
	EligibleForMeals     bool
	SealOfApproval       bool
}

func (q *Queries) RecipeSearch(ctx context.Context, db DBTX, arg *RecipeSearchParams) ([]*RecipeSearchRow, error) {
	rows, err := db.QueryContext(ctx, recipeSearch,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.Query,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*RecipeSearchRow{}
	for rows.Next() {
		var i RecipeSearchRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Source,
			&i.Description,
			&i.InspiredByRecipeID,
			&i.MinEstimatedPortions,
			&i.MaxEstimatedPortions,
			&i.PortionName,
			&i.PluralPortionName,
			&i.SealOfApproval,
			&i.EligibleForMeals,
			&i.YieldsComponentType,
			&i.LastIndexedAt,
			&i.LastValidatedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.CreatedByUser,
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

const updateRecipe = `-- name: UpdateRecipe :execrows
UPDATE recipes SET
	name = $1,
	slug = $2,
	source = $3,
	description = $4,
	inspired_by_recipe_id = $5,
	min_estimated_portions = $6,
	max_estimated_portions = $7,
	portion_name = $8,
	plural_portion_name = $9,
	seal_of_approval = $10,
	eligible_for_meals = $11,
	yields_component_type = $12,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND created_by_user = $13
	AND id = $14
`

type UpdateRecipeParams struct {
	YieldsComponentType  ComponentType
	Slug                 string
	Source               string
	Description          string
	ID                   string
	MinEstimatedPortions string
	Name                 string
	PortionName          string
	PluralPortionName    string
	CreatedByUser        string
	MaxEstimatedPortions sql.NullString
	InspiredByRecipeID   sql.NullString
	EligibleForMeals     bool
	SealOfApproval       bool
}

func (q *Queries) UpdateRecipe(ctx context.Context, db DBTX, arg *UpdateRecipeParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipe,
		arg.Name,
		arg.Slug,
		arg.Source,
		arg.Description,
		arg.InspiredByRecipeID,
		arg.MinEstimatedPortions,
		arg.MaxEstimatedPortions,
		arg.PortionName,
		arg.PluralPortionName,
		arg.SealOfApproval,
		arg.EligibleForMeals,
		arg.YieldsComponentType,
		arg.CreatedByUser,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateRecipeLastIndexedAt = `-- name: UpdateRecipeLastIndexedAt :execrows
UPDATE recipes SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateRecipeLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipeLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
