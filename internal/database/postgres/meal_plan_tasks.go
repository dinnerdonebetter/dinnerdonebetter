package postgres

import (
	"context"
	_ "embed"
	"sort"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealPlanTaskDataManager = (*Querier)(nil)
)

// scanMealPlanTaskWithRecipePrepTaskSteps takes a database Scanner (i.e. *sql.Row) and scans the result into a meal struct.
func (q *Querier) scanMealPlanTaskWithRecipePrepTaskSteps(ctx context.Context, rows database.ResultIterator) (x *types.MealPlanTask, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.MealPlanTask{}

	for rows.Next() {
		recipePrepTaskStep := &types.RecipePrepTaskStep{}

		targetVars := []any{
			&x.ID,
			&x.MealPlanOption.ID,
			&x.MealPlanOption.AssignedCook,
			&x.MealPlanOption.AssignedDishwasher,
			&x.MealPlanOption.Chosen,
			&x.MealPlanOption.TieBroken,
			&x.MealPlanOption.MealScale,
			&x.MealPlanOption.Meal.ID,
			&x.MealPlanOption.Notes,
			&x.MealPlanOption.CreatedAt,
			&x.MealPlanOption.LastUpdatedAt,
			&x.MealPlanOption.ArchivedAt,
			&x.MealPlanOption.BelongsToMealPlanEvent,
			&x.RecipePrepTask.ID,
			&x.RecipePrepTask.Name,
			&x.RecipePrepTask.Description,
			&x.RecipePrepTask.Notes,
			&x.RecipePrepTask.Optional,
			&x.RecipePrepTask.ExplicitStorageInstructions,
			&x.RecipePrepTask.MinimumTimeBufferBeforeRecipeInSeconds,
			&x.RecipePrepTask.MaximumTimeBufferBeforeRecipeInSeconds,
			&x.RecipePrepTask.StorageType,
			&x.RecipePrepTask.MinimumStorageTemperatureInCelsius,
			&x.RecipePrepTask.MaximumStorageTemperatureInCelsius,
			&x.RecipePrepTask.BelongsToRecipe,
			&x.RecipePrepTask.CreatedAt,
			&x.RecipePrepTask.LastUpdatedAt,
			&x.RecipePrepTask.ArchivedAt,
			&recipePrepTaskStep.ID,
			&recipePrepTaskStep.BelongsToRecipeStep,
			&recipePrepTaskStep.BelongsToRecipePrepTask,
			&recipePrepTaskStep.SatisfiesRecipeStep,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.CompletedAt,
			&x.Status,
			&x.CreationExplanation,
			&x.StatusExplanation,
			&x.AssignedToUser,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning complete recipe prep task step")
		}

		x.RecipePrepTask.TaskSteps = append(x.RecipePrepTask.TaskSteps, recipePrepTaskStep)
	}

	return x, nil
}

// scanMealPlanTasksForMealPlan takes a database Scanner (i.e. *sql.Row) and scans the result into a meal plan task struct.
func (q *Querier) scanMealPlanTasksForMealPlan(ctx context.Context, rows database.ResultIterator) ([]*types.MealPlanTask, error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	mealPlanTaskResults := []*types.MealPlanTask{}

	for rows.Next() {
		x := &types.MealPlanTask{}
		y := &types.RecipePrepTaskStep{}

		targetVars := []any{
			&x.ID,
			&x.MealPlanOption.ID,
			&x.MealPlanOption.AssignedCook,
			&x.MealPlanOption.AssignedDishwasher,
			&x.MealPlanOption.Chosen,
			&x.MealPlanOption.TieBroken,
			&x.MealPlanOption.MealScale,
			&x.MealPlanOption.Meal.ID,
			&x.MealPlanOption.Notes,
			&x.MealPlanOption.CreatedAt,
			&x.MealPlanOption.LastUpdatedAt,
			&x.MealPlanOption.ArchivedAt,
			&x.MealPlanOption.BelongsToMealPlanEvent,
			&x.RecipePrepTask.ID,
			&x.RecipePrepTask.Name,
			&x.RecipePrepTask.Description,
			&x.RecipePrepTask.Notes,
			&x.RecipePrepTask.Optional,
			&x.RecipePrepTask.ExplicitStorageInstructions,
			&x.RecipePrepTask.MinimumTimeBufferBeforeRecipeInSeconds,
			&x.RecipePrepTask.MaximumTimeBufferBeforeRecipeInSeconds,
			&x.RecipePrepTask.StorageType,
			&x.RecipePrepTask.MinimumStorageTemperatureInCelsius,
			&x.RecipePrepTask.MaximumStorageTemperatureInCelsius,
			&x.RecipePrepTask.BelongsToRecipe,
			&x.RecipePrepTask.CreatedAt,
			&x.RecipePrepTask.LastUpdatedAt,
			&x.RecipePrepTask.ArchivedAt,
			&y.ID,
			&y.BelongsToRecipeStep,
			&y.BelongsToRecipePrepTask,
			&y.SatisfiesRecipeStep,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.CompletedAt,
			&x.Status,
			&x.CreationExplanation,
			&x.StatusExplanation,
			&x.AssignedToUser,
		}

		if err := rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning meal plan task")
		}

		x.RecipePrepTask.TaskSteps = append(x.RecipePrepTask.TaskSteps, y)

		mealPlanTaskResults = append(mealPlanTaskResults, x)
	}

	// the TL;DR of this is that we get a list of every meal plan task step in a given meal plan,
	// for some unknown number of meal plan tasks. So we sort them by ID, congeal all the task
	// steps together and then return the results sorted by ID. There's probably a better and prettier
	// way to do this, but this is how we're doing it in this particular instance. At least I wrote this.
	mealPlanTaskMap := map[string]*types.MealPlanTask{}
	for _, mealPlanTask := range mealPlanTaskResults {
		if _, ok := mealPlanTaskMap[mealPlanTask.ID]; !ok {
			mealPlanTaskMap[mealPlanTask.ID] = mealPlanTask
		} else {
			mealPlanTaskMap[mealPlanTask.ID].RecipePrepTask.TaskSteps = append(mealPlanTaskMap[mealPlanTask.ID].RecipePrepTask.TaskSteps, mealPlanTask.RecipePrepTask.TaskSteps...)
		}
	}

	mealPlanTasks := types.MealPlanTaskList{}
	for _, mealPlanTask := range mealPlanTaskMap {
		mealPlanTasks = append(mealPlanTasks, mealPlanTask)
	}

	sort.Sort(mealPlanTasks)

	return mealPlanTasks, nil
}

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

	result, err := q.generatedQuerier.CheckMealPlanTaskExistence(ctx, q.db, &generated.CheckMealPlanTaskExistenceParams{
		MealPlanID:     mealPlanID,
		MealPlanTaskID: mealPlanTaskID,
	})
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

	args := []any{
		mealPlanTaskID,
	}

	rows, err := q.getRows(ctx, q.db, "meal plan task", getMealPlanTasksQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan task rows")
	}

	x, err = q.scanMealPlanTaskWithRecipePrepTaskSteps(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan task")
	}

	return x, nil
}

// createMealPlanTask creates a meal plan task.
func (q *Querier) createMealPlanTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	// create the meal plan task.
	if err := q.generatedQuerier.CreateMealPlanTask(ctx, querier, &generated.CreateMealPlanTaskParams{
		ID:                      input.ID,
		Status:                  types.MealPlanTaskStatusUnfinished,
		StatusExplanation:       input.StatusExplanation,
		CreationExplanation:     input.CreationExplanation,
		BelongsToMealPlanOption: input.MealPlanOptionID,
		BelongsToRecipePrepTask: input.RecipePrepTaskID,
		AssignedToUser:          nullStringFromStringPointer(input.AssignedToUser),
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	x := &types.MealPlanTask{
		CreatedAt:           q.timeFunc(),
		ID:                  input.ID,
		AssignedToUser:      input.AssignedToUser,
		Status:              types.MealPlanTaskStatusUnfinished,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOption: types.MealPlanOption{
			ID: input.MealPlanOptionID,
		},
		RecipePrepTask: types.RecipePrepTask{
			ID: input.RecipePrepTaskID,
		},
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
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("meal plan task created")

	return x, nil
}

//go:embed queries/meal_plan_tasks/list_all_by_meal_plan.sql
var listMealPlanTasksForMealPlanQuery string

// GetMealPlanTasksForMealPlan fetches a list of meal plan tasks.
func (q *Querier) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	args := []any{
		mealPlanID,
	}

	rows, getRowsErr := q.getRows(ctx, q.db, "meal plan tasks list", listMealPlanTasksForMealPlanQuery, args)
	if getRowsErr != nil {
		return nil, observability.PrepareAndLogError(getRowsErr, logger, span, "executing meal plan tasks list retrieval query")
	}

	x, scanErr := q.scanMealPlanTasksForMealPlan(ctx, rows)
	if scanErr != nil {
		return nil, observability.PrepareAndLogError(scanErr, logger, span, "scanning meal plan tasks")
	}

	logger.Info("meal plan tasks retrieved")

	return x, nil
}

// CreateMealPlanTasksForMealPlanOption creates meal plan tasks.
func (q *Querier) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
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
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(createMealPlanTaskErr, logger, span, "creating meal plan task")
		}

		outputs = append(outputs, mealPlanTask)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	logger.Info("meal plan tasks created")

	return outputs, nil
}

// MarkMealPlanAsHavingTasksCreated marks a meal plan as having all its tasks created.
func (q *Querier) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsPrepTasksCreated(ctx, q.db, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	logger.Info("meal plan tasks created")

	return nil
}

// MarkMealPlanAsHavingGroceryListInitialized marks a meal plan as having all its tasks created.
func (q *Querier) MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsGroceryListInitialized(ctx, q.db, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	logger.Info("meal plan marked as grocery list initialized")

	return nil
}

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

	var newStatus string
	if input.Status != nil {
		newStatus = *input.Status
	}

	if err := q.generatedQuerier.ChangeMealPlanTaskStatus(ctx, q.db, &generated.ChangeMealPlanTaskStatusParams{
		ID:                input.ID,
		Status:            generated.PrepStepStatus(newStatus),
		StatusExplanation: input.StatusExplanation,
		CompletedAt:       nullTimeFromTimePointer(settledAt),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing meal plan task status")
	}

	logger.Info("meal plan task status changed")

	return nil
}
