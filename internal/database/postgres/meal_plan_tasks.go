package postgres

import (
	"context"
	_ "embed"
	"time"

	"github.com/lib/pq"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.MealPlanTaskDataManager = (*Querier)(nil)
)

// scanMealPlanTask takes a database Scanner (i.e. *sql.Row) and scans the result into a valid instrument struct.
func (q *Querier) scanMealPlanTask(ctx context.Context, scan database.Scanner) (x *types.MealPlanTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanTask{}

	targetVars := []interface{}{
		&x.ID,
		&x.MealPlanOption.ID,
		&x.MealPlanOption.AssignedCook,
		&x.MealPlanOption.AssignedDishwasher,
		&x.MealPlanOption.Chosen,
		&x.MealPlanOption.TieBroken,
		&x.MealPlanOption.Meal.ID,
		&x.MealPlanOption.Notes,
		&x.MealPlanOption.PrepStepsCreated,
		&x.MealPlanOption.CreatedAt,
		&x.MealPlanOption.LastUpdatedAt,
		&x.MealPlanOption.ArchivedAt,
		&x.MealPlanOption.BelongsToMealPlanEvent,
		&x.RecipeStep.ID,
		&x.RecipeStep.Index,
		&x.RecipeStep.Preparation.ID,
		&x.RecipeStep.Preparation.Name,
		&x.RecipeStep.Preparation.Description,
		&x.RecipeStep.Preparation.IconPath,
		&x.RecipeStep.Preparation.YieldsNothing,
		&x.RecipeStep.Preparation.RestrictToIngredients,
		&x.RecipeStep.Preparation.ZeroIngredientsAllowable,
		&x.RecipeStep.Preparation.PastTense,
		&x.RecipeStep.Preparation.CreatedAt,
		&x.RecipeStep.Preparation.LastUpdatedAt,
		&x.RecipeStep.Preparation.ArchivedAt,
		&x.RecipeStep.MinimumEstimatedTimeInSeconds,
		&x.RecipeStep.MaximumEstimatedTimeInSeconds,
		&x.RecipeStep.MinimumTemperatureInCelsius,
		&x.RecipeStep.MaximumTemperatureInCelsius,
		&x.RecipeStep.Notes,
		&x.RecipeStep.ExplicitInstructions,
		&x.RecipeStep.Optional,
		&x.RecipeStep.CreatedAt,
		&x.RecipeStep.LastUpdatedAt,
		&x.RecipeStep.ArchivedAt,
		&x.RecipeStep.BelongsToRecipe,
		&x.AssignedToUser,
		&x.Status,
		&x.StatusExplanation,
		&x.CreationExplanation,
		&x.CannotCompleteBefore,
		&x.CannotCompleteAfter,
		&x.CreatedAt,
		&x.CompletedAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	return x, nil
}

// scanMealPlanTasks takes some database rows and turns them into a slice of advanced prep steps.
func (q *Querier) scanMealPlanTasks(ctx context.Context, rows database.ResultIterator) (validInstruments []*types.MealPlanTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, scanErr := q.scanMealPlanTask(ctx, rows)
		if scanErr != nil {
			return nil, scanErr
		}

		validInstruments = append(validInstruments, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, observability.PrepareError(err, span, "handling rows")
	}

	return validInstruments, nil
}

//go:embed queries/meal_plan_tasks/exists.sql
var mealPlanTasksExistsQuery string

// MealPlanTaskExists checks if an advanced prep step exists.
func (q *Querier) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if mealPlanTaskID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachMealPlanTaskIDToSpan(span, mealPlanTaskID)

	args := []interface{}{
		mealPlanID,
		mealPlanTaskID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, mealPlanTasksExistsQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing advanced step existence check")
	}

	logger.Info("advanced step existence retrieved")

	return result, nil
}

//go:embed queries/meal_plan_tasks/get_one.sql
var getMealPlanTasksQuery string

// GetMealPlanTask fetches an advanced prep step.
func (q *Querier) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (x *types.MealPlanTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachMealPlanTaskIDToSpan(span, mealPlanTaskID)

	args := []interface{}{
		mealPlanTaskID,
	}

	rows := q.getOneRow(ctx, q.db, "advanced prep step", getMealPlanTasksQuery, args)
	if x, err = q.scanMealPlanTask(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning advanced prep step")
	}

	logger.Info("advanced steps retrieved")

	return x, nil
}

//go:embed queries/meal_plan_tasks/list_all_by_meal_plan.sql
var listMealPlanTasksForMealPlanQuery string

// GetMealPlanTasksForMealPlan fetches a list of advanced prep steps.
func (q *Querier) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = []*types.MealPlanTask{}
	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	args := []interface{}{
		mealPlanID,
	}

	rows, getRowsErr := q.performReadQuery(ctx, q.db, "advanced prep steps list", listMealPlanTasksForMealPlanQuery, args)
	if getRowsErr != nil {
		return nil, observability.PrepareAndLogError(getRowsErr, logger, span, "executing advanced prep steps list retrieval query")
	}

	x, scanErr := q.scanMealPlanTasks(ctx, rows)
	if scanErr != nil {
		return nil, observability.PrepareAndLogError(scanErr, logger, span, "scanning advanced prep steps")
	}

	logger.Info("advanced steps retrieved")

	return x, nil
}

//go:embed queries/meal_plan_tasks/create.sql
var createMealPlanTaskQuery string

//go:embed queries/meal_plan_options/mark_as_steps_created.sql
var markMealPlanOptionAsHavingStepsCreatedQuery string

// CreateMealPlanTasksForMealPlanOption creates advanced prep steps.
func (q *Querier) CreateMealPlanTasksForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	outputs := []*types.MealPlanTask{}
	for _, input := range inputs {
		createMealPlanTaskArgs := []interface{}{
			input.ID,
			mealPlanOptionID,
			input.RecipeStepID,
			input.Status,
			input.StatusExplanation,
			input.CreationExplanation,
			pq.FormatTimestamp(input.CannotCompleteBefore.Truncate(time.Second)),
			pq.FormatTimestamp(input.CannotCompleteAfter.Truncate(time.Second)),
		}

		if err = q.performWriteQuery(ctx, tx, "create advanced prep step", createMealPlanTaskQuery, createMealPlanTaskArgs); err != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(err, logger, span, "create advanced prep step")
		}

		outputs = append(outputs, &types.MealPlanTask{
			ID:                   input.ID,
			CannotCompleteBefore: input.CannotCompleteBefore.Truncate(time.Second),
			CannotCompleteAfter:  input.CannotCompleteAfter.Truncate(time.Second),
			MealPlanOption:       types.MealPlanOption{ID: mealPlanOptionID},
			RecipeStep:           types.RecipeStep{ID: input.RecipeStepID},
			CreatedAt:            q.currentTime(),
			Status:               input.Status,
			StatusExplanation:    input.StatusExplanation,
			CreationExplanation:  input.CreationExplanation,
			CompletedAt:          input.CompletedAt,
		})
	}

	// mark prep steps as created for step
	markMealPlanOptionAsHavingStepsCreatedArgs := []interface{}{
		mealPlanOptionID,
	}

	if err = q.performWriteQuery(ctx, tx, "create advanced prep step", markMealPlanOptionAsHavingStepsCreatedQuery, markMealPlanOptionAsHavingStepsCreatedArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "create advanced prep step")
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	logger.Info("advanced steps created")

	return outputs, nil
}

//go:embed queries/meal_plan_tasks/change_status.sql
var changeMealPlanTaskStatusQuery string

// ChangeMealPlanTaskStatus changes an advanced prep step's status.
func (q *Querier) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}
	tracing.AttachMealPlanTaskIDToSpan(span, input.ID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	var settledAt *time.Time
	if input.Status == types.MealPlanTaskStatusFinished {
		t := q.timeFunc()
		settledAt = &t
	}

	changeMealPlanTaskStatusArgs := []interface{}{
		input.ID,
		input.Status,
		input.StatusExplanation,
		settledAt,
	}

	if err := q.performWriteQuery(ctx, q.db, "prep step status change", changeMealPlanTaskStatusQuery, changeMealPlanTaskStatusArgs); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing prep step status")
	}

	logger.Info("prep step status changed")

	return nil
}
