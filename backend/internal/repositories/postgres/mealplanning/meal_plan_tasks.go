package mealplanning

import (
	"context"
	"time"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	platformtypes "github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ types.MealPlanTaskDataManager = (*repository)(nil)
)

// MealPlanTaskExists checks if a meal plan task exists.
func (q *repository) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanTaskID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	result, err := q.generatedQuerier.CheckMealPlanTaskExistence(ctx, q.readDB, &generated.CheckMealPlanTaskExistenceParams{
		MealPlanID:     mealPlanID,
		MealPlanTaskID: mealPlanTaskID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan task existence check")
	}

	logger.Info("meal plan task existence retrieved")

	return result, nil
}

// GetMealPlanTask fetches a meal plan task.
func (q *repository) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanTaskID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	result, err := q.generatedQuerier.GetMealPlanTask(ctx, q.readDB, mealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting meal plan task")
	}

	mealPlanTask := &types.MealPlanTask{
		RecipePrepTask:      types.RecipePrepTask{},
		CreatedAt:           result.CreatedAt,
		LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
		CompletedAt:         database.TimePointerFromNullTime(result.CompletedAt),
		AssignedToUser:      database.StringPointerFromNullString(result.AssignedToUser),
		ID:                  result.ID,
		Status:              string(result.Status),
		CreationExplanation: result.CreationExplanation,
		StatusExplanation:   result.StatusExplanation,
		MealPlanOption: types.MealPlanOption{
			CreatedAt:              result.MealPlanOptionCreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.MealPlanOptionLastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.MealPlanOptionAssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.MealPlanOptionArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.MealPlanOptionAssignedDishwasher),
			Notes:                  result.MealPlanOptionNotes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.MealPlanOptionBelongsToMealPlanEvent),
			ID:                     result.MealPlanOptionID,
			Votes:                  nil,
			Meal: types.Meal{
				ID: result.MealPlanOptionMealID,
			},
			MealScale: database.Float32FromString(result.MealPlanOptionMealScale),
			Chosen:    result.MealPlanOptionChosen,
			TieBroken: result.MealPlanOptionTiebroken,
		},
	}

	return mealPlanTask, nil
}

// createMealPlanTask creates a meal plan task.
func (q *repository) createMealPlanTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
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
		AssignedToUser:          database.NullStringFromStringPointer(input.AssignedToUser),
	}); err != nil {
		q.RollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	x := &types.MealPlanTask{
		CreatedAt:           q.CurrentTime(),
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

	tracing.AttachToSpan(span, keys.MealPlanIDKey, x.ID)
	logger.Info("meal plan task created")

	return x, nil
}

// CreateMealPlanTask creates a meal plan task.
func (q *repository) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	x, err := q.createMealPlanTask(ctx, tx, input)
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("meal plan task created")

	return x, nil
}

// GetMealPlanTasksForMealPlan fetches a list of meal plan tasks.
func (q *repository) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanTask], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := q.generatedQuerier.ListAllMealPlanTasksByMealPlan(ctx, q.readDB, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan tasks list retrieval query")
	}

	// Group results by task ID (since each task can have multiple rows - one per prep task step)
	taskMap := make(map[string]*types.MealPlanTask)

	for _, result := range results {
		var mealPlanTask *types.MealPlanTask
		var exists bool

		if mealPlanTask, exists = taskMap[result.ID]; !exists {
			// Create new task
			mealPlanTask = &types.MealPlanTask{
				CreatedAt:           result.CreatedAt,
				LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
				CompletedAt:         database.TimePointerFromNullTime(result.CompletedAt),
				AssignedToUser:      database.StringPointerFromNullString(result.AssignedToUser),
				ID:                  result.ID,
				Status:              string(result.Status),
				CreationExplanation: result.CreationExplanation,
				StatusExplanation:   result.StatusExplanation,
				MealPlanOption: types.MealPlanOption{
					CreatedAt:              result.MealPlanOptionCreatedAt,
					LastUpdatedAt:          database.TimePointerFromNullTime(result.MealPlanOptionLastUpdatedAt),
					AssignedCook:           database.StringPointerFromNullString(result.MealPlanOptionAssignedCook),
					ArchivedAt:             database.TimePointerFromNullTime(result.MealPlanOptionArchivedAt),
					AssignedDishwasher:     database.StringPointerFromNullString(result.MealPlanOptionAssignedDishwasher),
					Notes:                  result.MealPlanOptionNotes,
					BelongsToMealPlanEvent: database.StringFromNullString(result.MealPlanOptionBelongsToMealPlanEvent),
					ID:                     result.MealPlanOptionID,
					Votes:                  nil,
					Meal: types.Meal{
						ID: result.MealPlanOptionMealID,
					},
					MealScale: database.Float32FromString(result.MealPlanOptionMealScale),
					Chosen:    result.MealPlanOptionChosen,
					TieBroken: result.MealPlanOptionTiebroken,
				},
				RecipePrepTask: types.RecipePrepTask{
					ID:                          result.PrepTaskID,
					BelongsToRecipe:             result.PrepTaskBelongsToRecipe,
					Name:                        result.PrepTaskName,
					Description:                 result.PrepTaskDescription,
					Notes:                       result.PrepTaskNotes,
					ExplicitStorageInstructions: result.PrepTaskExplicitStorageInstructions,
					Optional:                    result.PrepTaskOptional,
					CreatedAt:                   result.PrepTaskCreatedAt,
					LastUpdatedAt:               database.TimePointerFromNullTime(result.PrepTaskLastUpdatedAt),
					ArchivedAt:                  database.TimePointerFromNullTime(result.PrepTaskArchivedAt),
					StorageTemperatureInCelsius: platformtypes.OptionalFloat32Range{
						Max: database.Float32PointerFromNullString(result.PrepTaskMaximumStorageTemperatureInCelsius),
						Min: database.Float32PointerFromNullString(result.PrepTaskMinimumStorageTemperatureInCelsius),
					},
					TimeBufferBeforeRecipeInSeconds: platformtypes.Uint32RangeWithOptionalMax{
						Max: database.Uint32PointerFromNullInt32(result.PrepTaskMaximumTimeBufferBeforeRecipeInSeconds),
						Min: uint32(result.PrepTaskMinimumTimeBufferBeforeRecipeInSeconds),
					},
					TaskSteps: []*types.RecipePrepTaskStep{},
				},
			}

			// Set storage type if available
			if result.PrepTaskStorageType.Valid {
				mealPlanTask.RecipePrepTask.StorageType = string(result.PrepTaskStorageType.StorageContainerType)
			}

			taskMap[result.ID] = mealPlanTask
		}

		// Add prep task step if it exists and hasn't been added yet
		if result.PrepTaskStepID != "" {
			// Check if this step is already in the task steps
			stepExists := false
			for _, existingStep := range mealPlanTask.RecipePrepTask.TaskSteps {
				if existingStep.ID == result.PrepTaskStepID {
					stepExists = true
					break
				}
			}

			if !stepExists {
				mealPlanTask.RecipePrepTask.TaskSteps = append(mealPlanTask.RecipePrepTask.TaskSteps, &types.RecipePrepTaskStep{
					ID:                      result.PrepTaskStepID,
					BelongsToRecipeStep:     result.PrepTaskStepBelongsToRecipeStep,
					BelongsToRecipePrepTask: result.PrepTaskStepBelongsToRecipePrepTask,
					SatisfiesRecipeStep:     result.PrepTaskStepSatisfiesRecipeStep,
				})
			}
		}
	}

	// Convert map to slice
	data := make([]*types.MealPlanTask, 0, len(taskMap))
	for _, task := range taskMap {
		data = append(data, task)
	}

	var filteredCount, totalCount uint64

	x := filtering.NewQueryFilteredResult(data, filteredCount, totalCount, func(t *types.MealPlanTask) string {
		return t.ID
	}, filter)

	return x, nil
}

// CreateMealPlanTasksForMealPlanOption creates meal plan tasks.
func (q *repository) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	outputs := []*types.MealPlanTask{}
	for _, input := range inputs {
		mealPlanTask, createMealPlanTaskErr := q.createMealPlanTask(ctx, tx, input)
		if createMealPlanTaskErr != nil {
			q.RollbackTransaction(ctx, tx)
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
func (q *repository) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsPrepTasksCreated(ctx, q.writeDB, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	return nil
}

// MarkMealPlanAsHavingGroceryListInitialized marks a meal plan as having all its tasks created.
func (q *repository) MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsGroceryListInitialized(ctx, q.writeDB, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	logger.Info("meal plan marked as grocery list initialized")

	return nil
}

// ChangeMealPlanTaskStatus changes a meal plan task's status.
func (q *repository) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, input.MealPlanTaskID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.MealPlanTaskID)

	var settledAt *time.Time
	if input.Status != nil && *input.Status == types.MealPlanTaskStatusFinished {
		settledAt = new(q.CurrentTime())
	}

	var newStatus string
	if input.Status != nil {
		newStatus = *input.Status
	}

	if err := q.generatedQuerier.ChangeMealPlanTaskStatus(ctx, q.writeDB, &generated.ChangeMealPlanTaskStatusParams{
		ID:                input.MealPlanTaskID,
		Status:            generated.PrepStepStatus(newStatus),
		StatusExplanation: input.StatusExplanation,
		CompletedAt:       database.NullTimeFromTimePointer(settledAt),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing meal plan task status")
	}

	logger.Info("meal plan task status changed")

	return nil
}
