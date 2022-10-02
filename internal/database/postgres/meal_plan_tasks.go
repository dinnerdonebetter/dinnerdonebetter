package postgres

import (
	"context"
	_ "embed"
	"time"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.MealPlanTaskDataManager = (*Querier)(nil)
)

// scanMealPlanTaskWithRecipes takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMealPlanTaskWithRecipes(ctx context.Context, rows database.ResultIterator) (x *types.MealPlanTask, recipeIDs []string, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanTask{}

	for rows.Next() {
		var recipeID string

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
			&x.CannotCompleteBefore,
			&x.CannotCompleteAfter,
			&x.CreatedAt,
			&x.CompletedAt,
			&x.Status,
			&x.CreationExplanation,
			&x.StatusExplanation,
			&x.AssignedToUser,
			&recipeID,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, nil, observability.PrepareError(err, span, "scanning complete meal")
		}

		recipeIDs = append(recipeIDs, recipeID)
	}

	return x, recipeIDs, nil
}

// scanMealPlanTasksWithRecipes takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMealPlanTasksWithRecipes(ctx context.Context, rows database.ResultIterator) (mealPlanTasks []*types.MealPlanTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	var lastMealPlanTaskID string

	x := &types.MealPlanTask{}

	for rows.Next() {
		var recipeStepID string

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
			&x.CannotCompleteBefore,
			&x.CannotCompleteAfter,
			&x.CreatedAt,
			&x.CompletedAt,
			&x.Status,
			&x.CreationExplanation,
			&x.StatusExplanation,
			&x.AssignedToUser,
			&recipeStepID,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning complete meal")
		}

		x.RecipeSteps = append(x.RecipeSteps, &types.RecipeStep{ID: recipeStepID})

		if lastMealPlanTaskID == "" {
			lastMealPlanTaskID = x.ID
		}

		if len(mealPlanTasks) > 0 && lastMealPlanTaskID != mealPlanTasks[len(mealPlanTasks)-1].ID {
			mealPlanTasks = append(mealPlanTasks, x)
		}
	}

	mealPlanTasks = append(mealPlanTasks, x)

	return mealPlanTasks, nil
}

//go:embed queries/meal_plan_tasks/exists.sql
var mealPlanTasksExistsQuery string

// MealPlanTaskExists checks if a meal plan task exists.
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
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan task existence check")
	}

	logger.Info("meal plan task existence retrieved")

	return result, nil
}

//go:embed queries/meal_plan_tasks/get_one.sql
var getMealPlanTasksQuery string

// GetMealPlanTask fetches a meal plan task.
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

	rows, err := q.getRows(ctx, q.db, "meal plan task", getMealPlanTasksQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan task rows")
	}

	x, recipeStepIDs, err := q.scanMealPlanTaskWithRecipes(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning meal plan task")
	}

	for _, stepID := range recipeStepIDs {
		recipeStep, getRecipeStepErr := q.getRecipeStepByID(ctx, q.db, stepID)
		if getRecipeStepErr != nil {
			return nil, observability.PrepareAndLogError(getRecipeStepErr, logger, span, "getting recipe step for meal plan task")
		}

		x.RecipeSteps = append(x.RecipeSteps, recipeStep)
	}

	logger.Info("meal plan tasks retrieved")

	return x, nil
}

//go:embed queries/meal_plan_tasks/create.sql
var createMealPlanTaskQuery string

// createMealPlanTask creates a meal plan task.
func (q *Querier) createMealPlanTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	args := []interface{}{
		input.ID,
		types.MealPlanTaskStatusUnfinished,
		input.StatusExplanation,
		input.CreationExplanation,
		input.CannotCompleteBefore,
		input.CannotCompleteAfter,
		input.MealPlanOptionID,
		input.AssignedToUser,
	}

	// create the meal plan task.
	if err := q.performWriteQuery(ctx, querier, "meal plan task creation", createMealPlanTaskQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	x := &types.MealPlanTask{
		CannotCompleteBefore: input.CannotCompleteBefore,
		CannotCompleteAfter:  input.CannotCompleteAfter,
		CreatedAt:            q.timeFunc(),
		ID:                   input.ID,
		AssignedToUser:       input.AssignedToUser,
		Status:               types.MealPlanTaskStatusUnfinished,
		CreationExplanation:  input.CreationExplanation,
		StatusExplanation:    input.StatusExplanation,
	}

	for _, recipeStep := range input.RecipeSteps {
		if err := q.createMealPlanTaskRecipeStep(ctx, querier, recipeStep); err != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
		}
	}

	tracing.AttachMealPlanIDToSpan(span, x.ID)
	logger.Info("meal plan task created")

	return x, nil
}

// CreateMealPlanTask creates a meal plan task.
func (q *Querier) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	x, err := q.createMealPlanTask(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("meal plan task created")

	return x, nil
}

//go:embed queries/meal_plan_task_recipe_steps/create.sql
var createMealPlanTaskRecipeStepQuery string

// createMealPlanTaskRecipeStep creates a meal plan task recipe step.
func (q *Querier) createMealPlanTaskRecipeStep(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanTaskRecipeStepDatabaseCreationInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.BelongsToMealPlanTask,
		input.SatisfiesRecipeStep,
	}
	tracing.AttachMealPlanTaskIDToSpan(span, input.BelongsToMealPlanTask)

	// create the meal plan.
	if err := q.performWriteQuery(ctx, querier, "meal plan task recipe step creation", createMealPlanTaskRecipeStepQuery, args); err != nil {
		q.rollbackTransaction(ctx, querier)
		return observability.PrepareAndLogError(err, logger, span, "creating meal plan task recipe step")
	}

	logger.Info("meal plan task recipe step created")

	return nil
}

//go:embed queries/meal_plan_tasks/list_all_by_meal_plan.sql
var listMealPlanTasksForMealPlanQuery string

// GetMealPlanTasksForMealPlan fetches a list of meal plan tasks.
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

	rows, getRowsErr := q.getRows(ctx, q.db, "meal plan tasks list", listMealPlanTasksForMealPlanQuery, args)
	if getRowsErr != nil {
		return nil, observability.PrepareAndLogError(getRowsErr, logger, span, "executing meal plan tasks list retrieval query")
	}

	x, scanErr := q.scanMealPlanTasksWithRecipes(ctx, rows)
	if scanErr != nil {
		return nil, observability.PrepareAndLogError(scanErr, logger, span, "scanning meal plan tasks")
	}

	for i, mealPlanTask := range x {
		recipeSteps := []*types.RecipeStep{}
		for _, step := range mealPlanTask.RecipeSteps {
			recipeStep, recipeStepGetErr := q.getRecipeStepByID(ctx, q.db, step.ID)
			if recipeStepGetErr != nil {
				return nil, observability.PrepareAndLogError(recipeStepGetErr, logger, span, "getting recipe step for meal plan tasks")
			}

			recipeSteps = append(recipeSteps, recipeStep)
		}

		x[i].RecipeSteps = recipeSteps
	}

	logger.Info("meal plan tasks retrieved")

	return x, nil
}

//go:embed queries/meal_plan_options/mark_as_steps_created.sql
var markMealPlanOptionAsHavingStepsCreatedQuery string

// CreateMealPlanTasksForMealPlanOption creates meal plan tasks.
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
		mealPlanTask, createMealPlanTaskErr := q.createMealPlanTask(ctx, tx, input)
		if createMealPlanTaskErr != nil {
			return nil, observability.PrepareAndLogError(createMealPlanTaskErr, logger, span, "creating meal plan task")
		}

		outputs = append(outputs, mealPlanTask)
	}

	// mark prep steps as created for step
	markMealPlanOptionAsHavingStepsCreatedArgs := []interface{}{
		mealPlanOptionID,
	}

	if err = q.performWriteQuery(ctx, tx, "create meal plan task", markMealPlanOptionAsHavingStepsCreatedQuery, markMealPlanOptionAsHavingStepsCreatedArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "create meal plan task")
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	logger.Info("meal plan tasks created")

	return outputs, nil
}

//go:embed queries/meal_plan_tasks/change_status.sql
var changeMealPlanTaskStatusQuery string

// ChangeMealPlanTaskStatus changes a meal plan task's status.
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
	if input.Status != nil && *input.Status == types.MealPlanTaskStatusFinished {
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
