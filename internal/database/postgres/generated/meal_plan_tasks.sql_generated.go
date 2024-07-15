// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: meal_plan_tasks.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const changeMealPlanTaskStatus = `-- name: ChangeMealPlanTaskStatus :exec

UPDATE meal_plan_tasks SET
	completed_at = $1,
	status_explanation = $2,
	status = $3
WHERE id = $4
`

type ChangeMealPlanTaskStatusParams struct {
	CompletedAt       sql.NullTime
	StatusExplanation string
	Status            PrepStepStatus
	ID                string
}

func (q *Queries) ChangeMealPlanTaskStatus(ctx context.Context, db DBTX, arg *ChangeMealPlanTaskStatusParams) error {
	_, err := db.ExecContext(ctx, changeMealPlanTaskStatus,
		arg.CompletedAt,
		arg.StatusExplanation,
		arg.Status,
		arg.ID,
	)
	return err
}

const checkMealPlanTaskExistence = `-- name: CheckMealPlanTaskExistence :one

SELECT EXISTS (
	SELECT meal_plan_tasks.id
	FROM meal_plan_tasks
		FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
		FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_tasks.completed_at IS NULL
		AND meal_plans.id = $1
		AND meal_plans.archived_at IS NULL
		AND meal_plan_tasks.id = $2
)
`

type CheckMealPlanTaskExistenceParams struct {
	MealPlanID     string
	MealPlanTaskID string
}

func (q *Queries) CheckMealPlanTaskExistence(ctx context.Context, db DBTX, arg *CheckMealPlanTaskExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanTaskExistence, arg.MealPlanID, arg.MealPlanTaskID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createMealPlanTask = `-- name: CreateMealPlanTask :exec

INSERT INTO meal_plan_tasks (
	id,
	status,
	status_explanation,
	creation_explanation,
	belongs_to_meal_plan_option,
	belongs_to_recipe_prep_task,
	assigned_to_user
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
`

type CreateMealPlanTaskParams struct {
	ID                      string
	Status                  PrepStepStatus
	StatusExplanation       string
	CreationExplanation     string
	BelongsToMealPlanOption string
	BelongsToRecipePrepTask string
	AssignedToUser          sql.NullString
}

func (q *Queries) CreateMealPlanTask(ctx context.Context, db DBTX, arg *CreateMealPlanTaskParams) error {
	_, err := db.ExecContext(ctx, createMealPlanTask,
		arg.ID,
		arg.Status,
		arg.StatusExplanation,
		arg.CreationExplanation,
		arg.BelongsToMealPlanOption,
		arg.BelongsToRecipePrepTask,
		arg.AssignedToUser,
	)
	return err
}

const getMealPlanTask = `-- name: GetMealPlanTask :one

SELECT
	meal_plan_tasks.id,
	meal_plan_options.id as meal_plan_option_id,
	meal_plan_options.assigned_cook as meal_plan_option_assigned_cook,
	meal_plan_options.assigned_dishwasher as meal_plan_option_assigned_dishwasher,
	meal_plan_options.chosen as meal_plan_option_chosen,
	meal_plan_options.tiebroken as meal_plan_option_tiebroken,
	meal_plan_options.meal_scale as meal_plan_option_meal_scale,
	meal_plan_options.meal_id as meal_plan_option_meal_id,
	meal_plan_options.notes as meal_plan_option_notes,
	meal_plan_options.created_at as meal_plan_option_created_at,
	meal_plan_options.last_updated_at as meal_plan_option_last_updated_at,
	meal_plan_options.archived_at as meal_plan_option_archived_at,
	meal_plan_options.belongs_to_meal_plan_event as meal_plan_option_belongs_to_meal_plan_event,
	recipe_prep_tasks.id as prep_task_id,
	recipe_prep_tasks.name as prep_task_name,
	recipe_prep_tasks.description as prep_task_description,
	recipe_prep_tasks.notes as prep_task_notes,
	recipe_prep_tasks.optional as prep_task_optional,
	recipe_prep_tasks.explicit_storage_instructions as prep_task_explicit_storage_instructions,
	recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds as prep_task_minimum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds as prep_task_maximum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.storage_type as prep_task_storage_type,
	recipe_prep_tasks.minimum_storage_temperature_in_celsius as prep_task_minimum_storage_temperature_in_celsius,
	recipe_prep_tasks.maximum_storage_temperature_in_celsius as prep_task_maximum_storage_temperature_in_celsius,
	recipe_prep_tasks.created_at as prep_task_created_at,
	recipe_prep_tasks.last_updated_at as prep_task_last_updated_at,
	recipe_prep_tasks.archived_at as prep_task_archived_at,
	recipe_prep_tasks.belongs_to_recipe as prep_task_belongs_to_recipe,
	recipe_prep_task_steps.id as prep_task_step_id,
	recipe_prep_task_steps.belongs_to_recipe_step as prep_task_step_belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task as prep_task_step_belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step as prep_task_step_satisfies_recipe_step,
	meal_plan_tasks.status,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.belongs_to_meal_plan_option,
	meal_plan_tasks.belongs_to_recipe_prep_task,
	meal_plan_tasks.completed_at,
	meal_plan_tasks.created_at,
	meal_plan_tasks.last_updated_at,
	meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	JOIN meals ON meal_plan_options.meal_id=meals.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plans.archived_at IS NULL
	AND meals.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND meal_plan_tasks.id = $1
`

type GetMealPlanTaskRow struct {
	MealPlanOptionCreatedAt                        time.Time
	CreatedAt                                      time.Time
	PrepTaskCreatedAt                              time.Time
	LastUpdatedAt                                  sql.NullTime
	CompletedAt                                    sql.NullTime
	PrepTaskArchivedAt                             sql.NullTime
	PrepTaskLastUpdatedAt                          sql.NullTime
	MealPlanOptionArchivedAt                       sql.NullTime
	MealPlanOptionLastUpdatedAt                    sql.NullTime
	PrepTaskStepBelongsToRecipePrepTask            string
	PrepTaskStepID                                 string
	MealPlanOptionMealID                           string
	BelongsToRecipePrepTask                        string
	PrepTaskID                                     string
	PrepTaskName                                   string
	PrepTaskDescription                            string
	PrepTaskNotes                                  string
	MealPlanOptionNotes                            string
	PrepTaskExplicitStorageInstructions            string
	BelongsToMealPlanOption                        string
	CreationExplanation                            string
	StatusExplanation                              string
	Status                                         PrepStepStatus
	ID                                             string
	MealPlanOptionMealScale                        string
	PrepTaskStepBelongsToRecipeStep                string
	MealPlanOptionID                               string
	PrepTaskBelongsToRecipe                        string
	MealPlanOptionAssignedDishwasher               sql.NullString
	AssignedToUser                                 sql.NullString
	PrepTaskMaximumStorageTemperatureInCelsius     sql.NullString
	MealPlanOptionAssignedCook                     sql.NullString
	PrepTaskMinimumStorageTemperatureInCelsius     sql.NullString
	PrepTaskStorageType                            NullStorageContainerType
	MealPlanOptionBelongsToMealPlanEvent           sql.NullString
	PrepTaskMaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	PrepTaskMinimumTimeBufferBeforeRecipeInSeconds int32
	PrepTaskOptional                               bool
	PrepTaskStepSatisfiesRecipeStep                bool
	MealPlanOptionChosen                           bool
	MealPlanOptionTiebroken                        bool
}

func (q *Queries) GetMealPlanTask(ctx context.Context, db DBTX, mealPlanTaskID string) (*GetMealPlanTaskRow, error) {
	row := db.QueryRowContext(ctx, getMealPlanTask, mealPlanTaskID)
	var i GetMealPlanTaskRow
	err := row.Scan(
		&i.ID,
		&i.MealPlanOptionID,
		&i.MealPlanOptionAssignedCook,
		&i.MealPlanOptionAssignedDishwasher,
		&i.MealPlanOptionChosen,
		&i.MealPlanOptionTiebroken,
		&i.MealPlanOptionMealScale,
		&i.MealPlanOptionMealID,
		&i.MealPlanOptionNotes,
		&i.MealPlanOptionCreatedAt,
		&i.MealPlanOptionLastUpdatedAt,
		&i.MealPlanOptionArchivedAt,
		&i.MealPlanOptionBelongsToMealPlanEvent,
		&i.PrepTaskID,
		&i.PrepTaskName,
		&i.PrepTaskDescription,
		&i.PrepTaskNotes,
		&i.PrepTaskOptional,
		&i.PrepTaskExplicitStorageInstructions,
		&i.PrepTaskMinimumTimeBufferBeforeRecipeInSeconds,
		&i.PrepTaskMaximumTimeBufferBeforeRecipeInSeconds,
		&i.PrepTaskStorageType,
		&i.PrepTaskMinimumStorageTemperatureInCelsius,
		&i.PrepTaskMaximumStorageTemperatureInCelsius,
		&i.PrepTaskCreatedAt,
		&i.PrepTaskLastUpdatedAt,
		&i.PrepTaskArchivedAt,
		&i.PrepTaskBelongsToRecipe,
		&i.PrepTaskStepID,
		&i.PrepTaskStepBelongsToRecipeStep,
		&i.PrepTaskStepBelongsToRecipePrepTask,
		&i.PrepTaskStepSatisfiesRecipeStep,
		&i.Status,
		&i.StatusExplanation,
		&i.CreationExplanation,
		&i.BelongsToMealPlanOption,
		&i.BelongsToRecipePrepTask,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.AssignedToUser,
	)
	return &i, err
}

const listAllMealPlanTasksByMealPlan = `-- name: ListAllMealPlanTasksByMealPlan :many

SELECT
	meal_plan_tasks.id,
	meal_plan_options.id as meal_plan_option_id,
	meal_plan_options.assigned_cook as meal_plan_option_assigned_cook,
	meal_plan_options.assigned_dishwasher as meal_plan_option_assigned_dishwasher,
	meal_plan_options.chosen as meal_plan_option_chosen,
	meal_plan_options.tiebroken as meal_plan_option_tiebroken,
	meal_plan_options.meal_scale as meal_plan_option_meal_scale,
	meal_plan_options.meal_id as meal_plan_option_meal_id,
	meal_plan_options.notes as meal_plan_option_notes,
	meal_plan_options.created_at as meal_plan_option_created_at,
	meal_plan_options.last_updated_at as meal_plan_option_last_updated_at,
	meal_plan_options.archived_at as meal_plan_option_archived_at,
	meal_plan_options.belongs_to_meal_plan_event as meal_plan_option_belongs_to_meal_plan_event,
	recipe_prep_tasks.id as prep_task_id,
	recipe_prep_tasks.name as prep_task_name,
	recipe_prep_tasks.description as prep_task_description,
	recipe_prep_tasks.notes as prep_task_notes,
	recipe_prep_tasks.optional as prep_task_optional,
	recipe_prep_tasks.explicit_storage_instructions as prep_task_explicit_storage_instructions,
	recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds as prep_task_minimum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds as prep_task_maximum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.storage_type as prep_task_storage_type,
	recipe_prep_tasks.minimum_storage_temperature_in_celsius as prep_task_minimum_storage_temperature_in_celsius,
	recipe_prep_tasks.maximum_storage_temperature_in_celsius as prep_task_maximum_storage_temperature_in_celsius,
	recipe_prep_tasks.created_at as prep_task_created_at,
	recipe_prep_tasks.last_updated_at as prep_task_last_updated_at,
	recipe_prep_tasks.archived_at as prep_task_archived_at,
	recipe_prep_tasks.belongs_to_recipe as prep_task_belongs_to_recipe,
	recipe_prep_task_steps.id as prep_task_step_id,
	recipe_prep_task_steps.belongs_to_recipe_step as prep_task_step_belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task as prep_task_step_belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step as prep_task_step_satisfies_recipe_step,
	meal_plan_tasks.status,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.belongs_to_meal_plan_option,
	meal_plan_tasks.belongs_to_recipe_prep_task,
	meal_plan_tasks.completed_at,
	meal_plan_tasks.created_at,
	meal_plan_tasks.last_updated_at,
	meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	JOIN meals ON meal_plan_options.meal_id=meals.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plans.archived_at IS NULL
	AND meals.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND meal_plans.id = $1
`

type ListAllMealPlanTasksByMealPlanRow struct {
	MealPlanOptionCreatedAt                        time.Time
	CreatedAt                                      time.Time
	PrepTaskCreatedAt                              time.Time
	LastUpdatedAt                                  sql.NullTime
	CompletedAt                                    sql.NullTime
	PrepTaskArchivedAt                             sql.NullTime
	PrepTaskLastUpdatedAt                          sql.NullTime
	MealPlanOptionArchivedAt                       sql.NullTime
	MealPlanOptionLastUpdatedAt                    sql.NullTime
	PrepTaskStepBelongsToRecipePrepTask            string
	PrepTaskStepID                                 string
	MealPlanOptionMealID                           string
	BelongsToRecipePrepTask                        string
	PrepTaskID                                     string
	PrepTaskName                                   string
	PrepTaskDescription                            string
	PrepTaskNotes                                  string
	MealPlanOptionNotes                            string
	PrepTaskExplicitStorageInstructions            string
	BelongsToMealPlanOption                        string
	CreationExplanation                            string
	StatusExplanation                              string
	Status                                         PrepStepStatus
	ID                                             string
	MealPlanOptionMealScale                        string
	PrepTaskStepBelongsToRecipeStep                string
	MealPlanOptionID                               string
	PrepTaskBelongsToRecipe                        string
	MealPlanOptionAssignedDishwasher               sql.NullString
	AssignedToUser                                 sql.NullString
	PrepTaskMaximumStorageTemperatureInCelsius     sql.NullString
	MealPlanOptionAssignedCook                     sql.NullString
	PrepTaskMinimumStorageTemperatureInCelsius     sql.NullString
	PrepTaskStorageType                            NullStorageContainerType
	MealPlanOptionBelongsToMealPlanEvent           sql.NullString
	PrepTaskMaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	PrepTaskMinimumTimeBufferBeforeRecipeInSeconds int32
	PrepTaskOptional                               bool
	PrepTaskStepSatisfiesRecipeStep                bool
	MealPlanOptionChosen                           bool
	MealPlanOptionTiebroken                        bool
}

func (q *Queries) ListAllMealPlanTasksByMealPlan(ctx context.Context, db DBTX, mealPlanID string) ([]*ListAllMealPlanTasksByMealPlanRow, error) {
	rows, err := db.QueryContext(ctx, listAllMealPlanTasksByMealPlan, mealPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListAllMealPlanTasksByMealPlanRow{}
	for rows.Next() {
		var i ListAllMealPlanTasksByMealPlanRow
		if err := rows.Scan(
			&i.ID,
			&i.MealPlanOptionID,
			&i.MealPlanOptionAssignedCook,
			&i.MealPlanOptionAssignedDishwasher,
			&i.MealPlanOptionChosen,
			&i.MealPlanOptionTiebroken,
			&i.MealPlanOptionMealScale,
			&i.MealPlanOptionMealID,
			&i.MealPlanOptionNotes,
			&i.MealPlanOptionCreatedAt,
			&i.MealPlanOptionLastUpdatedAt,
			&i.MealPlanOptionArchivedAt,
			&i.MealPlanOptionBelongsToMealPlanEvent,
			&i.PrepTaskID,
			&i.PrepTaskName,
			&i.PrepTaskDescription,
			&i.PrepTaskNotes,
			&i.PrepTaskOptional,
			&i.PrepTaskExplicitStorageInstructions,
			&i.PrepTaskMinimumTimeBufferBeforeRecipeInSeconds,
			&i.PrepTaskMaximumTimeBufferBeforeRecipeInSeconds,
			&i.PrepTaskStorageType,
			&i.PrepTaskMinimumStorageTemperatureInCelsius,
			&i.PrepTaskMaximumStorageTemperatureInCelsius,
			&i.PrepTaskCreatedAt,
			&i.PrepTaskLastUpdatedAt,
			&i.PrepTaskArchivedAt,
			&i.PrepTaskBelongsToRecipe,
			&i.PrepTaskStepID,
			&i.PrepTaskStepBelongsToRecipeStep,
			&i.PrepTaskStepBelongsToRecipePrepTask,
			&i.PrepTaskStepSatisfiesRecipeStep,
			&i.Status,
			&i.StatusExplanation,
			&i.CreationExplanation,
			&i.BelongsToMealPlanOption,
			&i.BelongsToRecipePrepTask,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.AssignedToUser,
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

const listIncompleteMealPlanTasksByMealPlanOption = `-- name: ListIncompleteMealPlanTasksByMealPlanOption :many

SELECT
	meal_plan_tasks.id,
	meal_plan_options.id as meal_plan_option_id,
	meal_plan_options.assigned_cook as meal_plan_option_assigned_cook,
	meal_plan_options.assigned_dishwasher as meal_plan_option_assigned_dishwasher,
	meal_plan_options.chosen as meal_plan_option_chosen,
	meal_plan_options.tiebroken as meal_plan_option_tiebroken,
	meal_plan_options.meal_scale as meal_plan_option_meal_scale,
	meal_plan_options.meal_id as meal_plan_option_meal_id,
	meal_plan_options.notes as meal_plan_option_notes,
	meal_plan_options.created_at as meal_plan_option_created_at,
	meal_plan_options.last_updated_at as meal_plan_option_last_updated_at,
	meal_plan_options.archived_at as meal_plan_option_archived_at,
	meal_plan_options.belongs_to_meal_plan_event as meal_plan_option_belongs_to_meal_plan_event,
	recipe_steps.id as recipe_step_id,
	recipe_steps.index as recipe_step_index,
	valid_preparations.id as valid_preparation_id,
	valid_preparations.name as valid_preparation_name,
	valid_preparations.description as valid_preparation_description,
	valid_preparations.icon_path as valid_preparation_icon_path,
	valid_preparations.yields_nothing as valid_preparation_yields_nothing,
	valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
	valid_preparations.past_tense as valid_preparation_past_tense,
	valid_preparations.slug as valid_preparation_slug,
	valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
	valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
	valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
	valid_preparations.temperature_required as valid_preparation_temperature_required,
	valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
	valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
	valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
	valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
	valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
	valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
	valid_preparations.last_indexed_at as valid_preparation_last_indexed_at,
	valid_preparations.created_at as valid_preparation_created_at,
	valid_preparations.last_updated_at as valid_preparation_last_updated_at,
	valid_preparations.archived_at as valid_preparation_archived_at,
	recipe_steps.preparation_id as recipe_step_preparation_id,
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
	recipe_steps.belongs_to_recipe as recipe_step_belongs_to_recipe,
	meal_plan_tasks.status,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.belongs_to_meal_plan_option,
	meal_plan_tasks.belongs_to_recipe_prep_task,
	meal_plan_tasks.completed_at,
	meal_plan_tasks.created_at,
	meal_plan_tasks.last_updated_at,
	meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	 FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	 FULL OUTER JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
	 FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
	 JOIN recipe_steps ON meal_plan_tasks.satisfies_recipe_step=recipe_steps.id
	 JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE meal_plan_tasks.belongs_to_meal_plan_option = $1
AND meal_plan_tasks.completed_at IS NULL
`

type ListIncompleteMealPlanTasksByMealPlanOptionRow struct {
	RecipeStepCreatedAt                         time.Time
	ValidPreparationCreatedAt                   time.Time
	CompletedAt                                 sql.NullTime
	ValidPreparationLastIndexedAt               sql.NullTime
	RecipeStepArchivedAt                        sql.NullTime
	RecipeStepLastUpdatedAt                     sql.NullTime
	CreatedAt                                   sql.NullTime
	ValidPreparationArchivedAt                  sql.NullTime
	ValidPreparationLastUpdatedAt               sql.NullTime
	MealPlanOptionCreatedAt                     sql.NullTime
	MealPlanOptionLastUpdatedAt                 sql.NullTime
	MealPlanOptionArchivedAt                    sql.NullTime
	LastUpdatedAt                               sql.NullTime
	RecipeStepConditionExpression               string
	RecipeStepPreparationID                     string
	ValidPreparationID                          string
	ValidPreparationName                        string
	ValidPreparationDescription                 string
	ValidPreparationIconPath                    string
	RecipeStepExplicitInstructions              string
	RecipeStepNotes                             string
	ValidPreparationPastTense                   string
	ValidPreparationSlug                        string
	RecipeStepID                                string
	RecipeStepBelongsToRecipe                   string
	MealPlanOptionBelongsToMealPlanEvent        sql.NullString
	CreationExplanation                         sql.NullString
	MealPlanOptionAssignedDishwasher            sql.NullString
	MealPlanOptionMealID                        sql.NullString
	MealPlanOptionNotes                         sql.NullString
	StatusExplanation                           sql.NullString
	Status                                      NullPrepStepStatus
	ID                                          sql.NullString
	AssignedToUser                              sql.NullString
	MealPlanOptionAssignedCook                  sql.NullString
	BelongsToRecipePrepTask                     sql.NullString
	MealPlanOptionID                            sql.NullString
	BelongsToMealPlanOption                     sql.NullString
	MealPlanOptionMealScale                     sql.NullString
	RecipeStepMaximumTemperatureInCelsius       sql.NullString
	RecipeStepMinimumTemperatureInCelsius       sql.NullString
	RecipeStepMaximumEstimatedTimeInSeconds     sql.NullInt64
	RecipeStepMinimumEstimatedTimeInSeconds     sql.NullInt64
	ValidPreparationMaximumInstrumentCount      sql.NullInt32
	ValidPreparationMaximumVesselCount          sql.NullInt32
	ValidPreparationMaximumIngredientCount      sql.NullInt32
	ValidPreparationMinimumInstrumentCount      int32
	RecipeStepIndex                             int32
	ValidPreparationMinimumIngredientCount      int32
	ValidPreparationMinimumVesselCount          int32
	MealPlanOptionTiebroken                     sql.NullBool
	MealPlanOptionChosen                        sql.NullBool
	ValidPreparationConsumesVessel              bool
	RecipeStepStartTimerAutomatically           bool
	ValidPreparationConditionExpressionRequired bool
	ValidPreparationTimeEstimateRequired        bool
	ValidPreparationTemperatureRequired         bool
	ValidPreparationOnlyForVessels              bool
	RecipeStepOptional                          bool
	ValidPreparationRestrictToIngredients       bool
	ValidPreparationYieldsNothing               bool
}

func (q *Queries) ListIncompleteMealPlanTasksByMealPlanOption(ctx context.Context, db DBTX, belongsToMealPlanOption string) ([]*ListIncompleteMealPlanTasksByMealPlanOptionRow, error) {
	rows, err := db.QueryContext(ctx, listIncompleteMealPlanTasksByMealPlanOption, belongsToMealPlanOption)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListIncompleteMealPlanTasksByMealPlanOptionRow{}
	for rows.Next() {
		var i ListIncompleteMealPlanTasksByMealPlanOptionRow
		if err := rows.Scan(
			&i.ID,
			&i.MealPlanOptionID,
			&i.MealPlanOptionAssignedCook,
			&i.MealPlanOptionAssignedDishwasher,
			&i.MealPlanOptionChosen,
			&i.MealPlanOptionTiebroken,
			&i.MealPlanOptionMealScale,
			&i.MealPlanOptionMealID,
			&i.MealPlanOptionNotes,
			&i.MealPlanOptionCreatedAt,
			&i.MealPlanOptionLastUpdatedAt,
			&i.MealPlanOptionArchivedAt,
			&i.MealPlanOptionBelongsToMealPlanEvent,
			&i.RecipeStepID,
			&i.RecipeStepIndex,
			&i.ValidPreparationID,
			&i.ValidPreparationName,
			&i.ValidPreparationDescription,
			&i.ValidPreparationIconPath,
			&i.ValidPreparationYieldsNothing,
			&i.ValidPreparationRestrictToIngredients,
			&i.ValidPreparationPastTense,
			&i.ValidPreparationSlug,
			&i.ValidPreparationMinimumIngredientCount,
			&i.ValidPreparationMaximumIngredientCount,
			&i.ValidPreparationMinimumInstrumentCount,
			&i.ValidPreparationMaximumInstrumentCount,
			&i.ValidPreparationTemperatureRequired,
			&i.ValidPreparationTimeEstimateRequired,
			&i.ValidPreparationConditionExpressionRequired,
			&i.ValidPreparationConsumesVessel,
			&i.ValidPreparationOnlyForVessels,
			&i.ValidPreparationMinimumVesselCount,
			&i.ValidPreparationMaximumVesselCount,
			&i.ValidPreparationLastIndexedAt,
			&i.ValidPreparationCreatedAt,
			&i.ValidPreparationLastUpdatedAt,
			&i.ValidPreparationArchivedAt,
			&i.RecipeStepPreparationID,
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
			&i.Status,
			&i.StatusExplanation,
			&i.CreationExplanation,
			&i.BelongsToMealPlanOption,
			&i.BelongsToRecipePrepTask,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.AssignedToUser,
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
