package postgres

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipePrepTaskDataManager = (*Querier)(nil)
)

// RecipePrepTaskExists checks if a recipe prep task exists.
func (q *Querier) RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	result, err := q.generatedQuerier.CheckRecipePrepTaskExistence(ctx, q.db, &generated.CheckRecipePrepTaskExistenceParams{
		BelongsToRecipe: recipeID,
		ID:              recipePrepTaskID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe prep task existence check")
	}

	logger.Info("recipe prep task existence retrieved")

	return result, nil
}

// GetRecipePrepTask fetches a recipe prep task.
func (q *Querier) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (x *types.RecipePrepTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	results, err := q.generatedQuerier.GetRecipePrepTask(ctx, q.db, recipePrepTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe prep task")
	}

	for _, result := range results {
		if x == nil {
			x = &types.RecipePrepTask{
				CreatedAt:                              result.CreatedAt,
				MaximumStorageTemperatureInCelsius:     float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
				ArchivedAt:                             timePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt:                          timePointerFromNullTime(result.LastUpdatedAt),
				MinimumStorageTemperatureInCelsius:     float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
				MaximumTimeBufferBeforeRecipeInSeconds: uint32PointerFromNullInt32(result.MaximumTimeBufferBeforeRecipeInSeconds),
				ID:                                     result.ID,
				BelongsToRecipe:                        result.BelongsToRecipe,
				ExplicitStorageInstructions:            result.ExplicitStorageInstructions,
				Notes:                                  result.Notes,
				Name:                                   result.Name,
				Description:                            result.Description,
				TaskSteps:                              []*types.RecipePrepTaskStep{},
				MinimumTimeBufferBeforeRecipeInSeconds: uint32(result.MinimumTimeBufferBeforeRecipeInSeconds),
				Optional:                               result.Optional,
			}

			logger.WithValue("storage_type", result.StorageType).Info("storage type")

			if result.StorageType.Valid {
				x.StorageType = string(result.StorageType.StorageContainerType)
			}
		}

		x.TaskSteps = append(x.TaskSteps, &types.RecipePrepTaskStep{
			ID:                      result.TaskStepID,
			BelongsToRecipeStep:     result.TaskStepBelongsToRecipeStep,
			BelongsToRecipePrepTask: result.TaskStepBelongsToRecipePrepTask,
			SatisfiesRecipeStep:     result.TaskStepSatisfiesRecipeStep,
		})
	}

	if x == nil || x.ID == "" {
		return nil, sql.ErrNoRows
	}

	logger.Info("recipe prep tasks retrieved")

	return x, nil
}

// createRecipePrepTask creates a recipe prep task.
func (q *Querier) createRecipePrepTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachMealPlanIDToSpan(span, input.ID)
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	// create the recipe prep task.
	if err := q.generatedQuerier.CreateRecipePrepTask(ctx, querier, &generated.CreateRecipePrepTaskParams{
		ID:                                     input.ID,
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		BelongsToRecipe:                        input.BelongsToRecipe,
		StorageType:                            generated.NullStorageContainerType{StorageContainerType: generated.StorageContainerType(input.StorageType), Valid: true},
		MinimumStorageTemperatureInCelsius:     nullStringFromFloat32Pointer(input.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius:     nullStringFromFloat32Pointer(input.MaximumStorageTemperatureInCelsius),
		MaximumTimeBufferBeforeRecipeInSeconds: nullInt32FromUint32Pointer(input.MaximumTimeBufferBeforeRecipeInSeconds),
		MinimumTimeBufferBeforeRecipeInSeconds: int32(input.MinimumTimeBufferBeforeRecipeInSeconds),
		Optional:                               input.Optional,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	x := &types.RecipePrepTask{
		CreatedAt:                              q.timeFunc(),
		ID:                                     input.ID,
		Name:                                   input.Name,
		Description:                            input.Description,
		Notes:                                  input.Notes,
		Optional:                               input.Optional,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		StorageType:                            input.StorageType,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        input.BelongsToRecipe,
	}

	for _, recipePrepTaskStep := range input.TaskSteps {
		s, err := q.createRecipePrepTaskStep(ctx, querier, recipePrepTaskStep)
		if err != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
		}

		x.TaskSteps = append(x.TaskSteps, s)
	}

	logger.Info("recipe prep task created")

	return x, nil
}

// CreateRecipePrepTask creates a recipe prep task.
func (q *Querier) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	x, err := q.createRecipePrepTask(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("recipe prep task created")

	return x, nil
}

// createRecipePrepTaskStep creates a recipe prep task step.
func (q *Querier) createRecipePrepTaskStep(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.RecipePrepTaskStepDatabaseCreationInput) (*types.RecipePrepTaskStep, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, input.ID)
	tracing.AttachRecipePrepTaskIDToSpan(span, input.BelongsToRecipePrepTask)

	// create the meal plan.
	if err := q.generatedQuerier.CreateRecipePrepTaskStep(ctx, querier, &generated.CreateRecipePrepTaskStepParams{
		ID:                      input.ID,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe prep task step")
	}

	x := &types.RecipePrepTaskStep{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}

	logger.Info("recipe prep task step created")

	return x, nil
}

// getRecipePrepTasksForRecipe gets a recipe prep task.
func (q *Querier) getRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*types.RecipePrepTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	results, err := q.generatedQuerier.ListAllRecipePrepTasksByRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe prep tasks list retrieval query")
	}

	x := []*types.RecipePrepTask{}

	var currentRecipePrepTask *types.RecipePrepTask
	for _, result := range results {
		prepTaskStep := &types.RecipePrepTaskStep{
			ID:                      result.TaskStepID,
			BelongsToRecipeStep:     result.TaskStepBelongsToRecipeStep,
			BelongsToRecipePrepTask: result.TaskStepBelongsToRecipePrepTask,
			SatisfiesRecipeStep:     result.TaskStepSatisfiesRecipeStep,
		}

		if currentRecipePrepTask != nil && currentRecipePrepTask.ID != result.ID {
			x = append(x, currentRecipePrepTask)
			currentRecipePrepTask = nil
		}

		if currentRecipePrepTask == nil {
			currentRecipePrepTask = &types.RecipePrepTask{
				CreatedAt:                              result.CreatedAt,
				MaximumStorageTemperatureInCelsius:     float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
				ArchivedAt:                             timePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt:                          timePointerFromNullTime(result.LastUpdatedAt),
				MinimumStorageTemperatureInCelsius:     float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
				MaximumTimeBufferBeforeRecipeInSeconds: uint32PointerFromNullInt32(result.MaximumTimeBufferBeforeRecipeInSeconds),
				ID:                                     result.ID,
				StorageType:                            string(result.StorageType.StorageContainerType),
				BelongsToRecipe:                        result.BelongsToRecipe,
				ExplicitStorageInstructions:            result.ExplicitStorageInstructions,
				Notes:                                  result.Notes,
				Name:                                   result.Name,
				Description:                            result.Description,
				TaskSteps:                              []*types.RecipePrepTaskStep{},
				MinimumTimeBufferBeforeRecipeInSeconds: uint32(result.MinimumTimeBufferBeforeRecipeInSeconds),
				Optional:                               result.Optional,
			}
		}
		currentRecipePrepTask.TaskSteps = append(currentRecipePrepTask.TaskSteps, prepTaskStep)
	}

	if currentRecipePrepTask != nil && currentRecipePrepTask.ID != "" {
		x = append(x, currentRecipePrepTask)
	}

	logger.Info("recipe prep tasks retrieved")

	return x, nil
}

// GetRecipePrepTasksForRecipe gets a recipe prep task.
func (q *Querier) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) (x []*types.RecipePrepTask, err error) {
	return q.getRecipePrepTasksForRecipe(ctx, recipeID)
}

// UpdateRecipePrepTask updates a recipe prep task.
func (q *Querier) UpdateRecipePrepTask(ctx context.Context, updated *types.RecipePrepTask) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	if updated == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateRecipePrepTask(ctx, q.db, &generated.UpdateRecipePrepTaskParams{
		Name:                                   updated.Name,
		Description:                            updated.Description,
		Notes:                                  updated.Notes,
		Optional:                               updated.Optional,
		ExplicitStorageInstructions:            updated.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: int32(updated.MinimumTimeBufferBeforeRecipeInSeconds),
		MaximumTimeBufferBeforeRecipeInSeconds: nullInt32FromUint32Pointer(updated.MaximumTimeBufferBeforeRecipeInSeconds),
		StorageType:                            generated.NullStorageContainerType{StorageContainerType: generated.StorageContainerType(updated.StorageType), Valid: true},
		MinimumStorageTemperatureInCelsius:     nullStringFromFloat32Pointer(updated.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius:     nullStringFromFloat32Pointer(updated.MaximumStorageTemperatureInCelsius),
		BelongsToRecipe:                        updated.BelongsToRecipe,
		ID:                                     updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	logger.Info("recipe prep task updated")

	return nil
}

// ArchiveRecipePrepTask marks a recipe prep task as archived.
func (q *Querier) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipePrepTaskID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipePrepTaskIDKey, recipePrepTaskID)
	tracing.AttachRecipePrepTaskIDToSpan(span, recipePrepTaskID)

	if _, err := q.generatedQuerier.ArchiveRecipePrepTask(ctx, q.db, recipePrepTaskID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe prep task")
	}

	return nil
}
